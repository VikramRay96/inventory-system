package impl

import (
	"context"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/common/pkg/utils"
	"inventory-system/inventory-service/internal/adapters/client"
	"inventory-system/inventory-service/internal/adapters/models"
	"inventory-system/inventory-service/internal/adapters/repository"
	commonDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	"inventory-system/inventory-service/internal/domain/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type InventoryService struct {
	InventoryRepository           repository.IInventoryRepository
	InventoryConfigurationService service.IInventoryConfigurationService
	S3                            client.IS3
}

type Resource struct {
	ServiceName string `json:"service_name"`
	FlowType    string `json:"flow_type"`
	TopicId     string `json:"topic_id"`
	TopicName   string `json:"topic_name"`
	Path        string `json:"path"`
	Id          string `json:"id"`
}

type ResourceModel struct {
	ServiceName string `json:"service_name" bson:"service_name"`
	FlowType    string `json:"flow_type" bson:"flow_type"`
	TopicId     string `json:"topic_id" bson:"topic_id"`
	TopicName   string `json:"topic_name" bson:"topic_name"`
	Path        string `json:"path" bson:"path"`
	Id          string `json:"id" bson:"id"`
	IsDeleted   bool   `json:"is_deleted" bson:"is_deleted"`
}

func NewInventoryService(inventoryRepository repository.IInventoryRepository, inventoryConfigurationService service.IInventoryConfigurationService, is3 client.IS3) *InventoryService {
	s := &InventoryService{
		InventoryRepository:           inventoryRepository,
		InventoryConfigurationService: inventoryConfigurationService,
		S3:                            is3,
	}
	return s
}
func (c InventoryService) CreateNewInventory(ctx context.Context, item interface{}, inventoryName string) *dto.ErrorResponseDto {
	methodName := "CreateNewInventory"
	log := logger.GetLogger()

	//check if inventory item deleted.
	_, errDto := c.InventoryConfigurationService.GetInventoryConfiguration(ctx, inventoryName)
	if errDto != nil {
		if errDto.StatusCode == status_code.IMS204 {
			log.Info("Inside "+methodName+" inventory item deleted :", inventoryName)
			return errDto
		}
		log.Info("Inside "+methodName+" unable to fetch inventory item for inventoryName :", inventoryName)
		return errDto
	}

	//Create New item
	errorDto := c.InventoryRepository.CreateNewInventoryGivenInventoryName(ctx, item, inventoryName)
	if errorDto != nil {
		log.Error("Inside "+methodName+" error occurred when trying to add new inventory: ", item)
		return errorDto
	}
	log.Info("Inside " + methodName + "successfully created new inventory inside " + inventoryName)
	return nil
}

func (c InventoryService) GetInventory(ctx context.Context, inventoryName string, filterMap map[string][]string) (bson.M, *dto.ErrorResponseDto) {
	methodName := "GetInventory"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto
	var InventoryIdentifier string
	var InventoryIdentifierValue string

	log.Info("Inside "+methodName+" getting inventory item for :", inventoryName)

	//Check if only 1 query param exists.
	if len(filterMap) > 1 {
		log.Error("Inside "+methodName+" multiple filters not allowed :", inventoryName, " and filterMap : ", filterMap)
		domainErr.SetError(status_code.IMS112)
		return nil, &domainErr
	}

	//Fetching inventory item to check if not deleted.
	inventoryConfiguration, errDto := c.InventoryConfigurationService.GetInventoryConfiguration(ctx, inventoryName)
	if errDto != nil {
		if errDto.StatusCode == status_code.IMS204 {
			log.Info("Inside "+methodName+" inventory item deleted :", inventoryName)
			return nil, errDto
		}
		log.Info("Inside "+methodName+" unable to fetch inventory item for inventoryName :", inventoryName)
		return nil, errDto
	}

	//Checking if filter attribute exist in identifier list and is unique
	for key, value := range filterMap {
		if UniqueKeyExists(inventoryConfiguration.InventoryIdentifiers, key) == false {
			log.Error("Inside "+methodName+" filter attribute not present in unique inventoryIdentifiers :", inventoryName, " and filterMap : ", filterMap)
			domainErr.SetError(status_code.IMS115)
			return nil, &domainErr
		}
		InventoryIdentifier = key
		InventoryIdentifierValue = value[0]
	}

	//Fetching item from inventory
	item, adapterError := c.InventoryRepository.FetchInventory(ctx, inventoryName, InventoryIdentifier, InventoryIdentifierValue)
	if adapterError != nil {
		log.Error("Inside "+methodName+" error while fetching item for :", inventoryName, " and filterMap : ", filterMap)
		return nil, adapterError
	}
	return item, nil

}
func (c InventoryService) GetInventoryV2(ctx *gin.Context, from string, to string, inventoryName string, filterMap map[string][]string) ([]bson.M, *commonDto.PaginationResponse, *dto.ErrorResponseDto) {
	methodName := "GetInventory"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto

	log.Info("Inside "+methodName+" getting inventory item for :", inventoryName)

	//Fetching inventory item to check if not deleted.
	inventoryConfiguration, errDto := c.InventoryConfigurationService.GetInventoryConfiguration(ctx, inventoryName)
	if errDto != nil {
		if errDto.StatusCode == status_code.IMS204 {
			log.Info("Inside "+methodName+" inventory item deleted :", inventoryName)
			return nil, nil, errDto
		}
		log.Info("Inside "+methodName+" unable to fetch inventory item for inventoryName :", inventoryName)
		return nil, nil, errDto
	}
	var pagination commonDto.Pagination
	if inventoryConfiguration.Pagination {
		pagination.Pagination = true
		page := ctx.Query("page")
		cPage, _ := strconv.ParseInt(page, 10, 64)

		pagination.PageNumber = cPage

		pageSize := ctx.Query("page_size")

		cPageSize, _ := strconv.ParseInt(pageSize, 10, 64)
		log.Info("cPageSize", cPageSize)
		if cPageSize == 0 {
			pagination.PageSize = 10
		} else {

			pagination.PageSize = cPageSize
		}
	}

	//Checking if filter attribute exist in identifier list
	for key, _ := range filterMap {
		log.Info("InventoryIdentifiers", inventoryConfiguration.InventoryIdentifiers)
		if !KeyExists(inventoryConfiguration.InventoryIdentifiers, key) {
			log.Error("Inside "+methodName+" filter attribute not present in inventoryIdentifiers :", inventoryName, " and filterMap : ", filterMap)
			domainErr.SetError(status_code.IMS113)
			return nil, nil, &domainErr
		}
	}

	//Fetching item from inventory
	item, paginationData, adapterError := c.InventoryRepository.FetchInventoryList(ctx, from, to, inventoryName, filterMap, pagination)
	if adapterError != nil {
		log.Error("Inside "+methodName+" error while fetching item for :", inventoryName, " and filterMap : ", filterMap)
		return nil, nil, adapterError
	}

	return item, paginationData, nil

}

func (c InventoryService) RemoveItemFromInventory(ctx context.Context, RemoveInventoryItemRequest *request_dto.RemoveInventoryItem, InventoryName string) *dto.ErrorResponseDto {
	methodName := "RemoveItemFromInventory"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto
	log.Info("Inside "+methodName+" getting inventory item for :", InventoryName)

	RemoveInventoryItemModel, BindingError := utils.TypeConverter[models.RemoveInventoryItem](RemoveInventoryItemRequest)

	if BindingError != nil {
		log.Info("Error while binding Remove Inventory Item Model", BindingError)
		domainErr.SetError(status_code.IMS400)
		return &domainErr
	}

	RemoveModelError := c.InventoryRepository.RemoveItemFromInventory(ctx, RemoveInventoryItemModel, InventoryName)

	if RemoveModelError != nil {
		log.Info("Error while removing Item from the inventory", RemoveModelError)
		return RemoveModelError
	}
	return nil
}

func (c InventoryService) UpdateInventory(Id string, InventoryName string, UpdateRequest *interface{}) *dto.ErrorResponseDto {
	log := logger.GetLogger()
	methodName := "UpdateInventory Repository"

	log.Info("Inside " + methodName)

	log.Info("Update Request Body", UpdateRequest)

	AdapterError := c.InventoryRepository.UpdateInventory(Id, InventoryName, UpdateRequest)
	if AdapterError != nil {
		return AdapterError
	}
	return nil
}

func (c InventoryService) RemoveSubjectTopicsByLessonNameAndSubjectId(ctx context.Context, RemoveSubjectRequest *request_dto.RemoveSubjectRequest, Type string) *dto.ErrorResponseDto {
	methodName := "RemoveSubjectTopicsByLessonNameAndSubjectId"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto
	log.Info("Inside " + methodName)

	RemoveInventorySubjectModel, BindingError := utils.TypeConverter[models.RemoveSubjectRequestModel](RemoveSubjectRequest)

	if BindingError != nil {
		log.Info("Error while binding Remove Inventory Subject Model", BindingError)
		domainErr.SetError(status_code.IMS400)
		return &domainErr
	}

	RemoveSubjectModelError := c.InventoryRepository.RemoveSubjectTopicsByLessonNameAndSubjectId(ctx, RemoveInventorySubjectModel, Type)

	if RemoveSubjectModelError != nil {
		log.Info("Error while removing Subject from the inventory", RemoveSubjectModelError)
		return RemoveSubjectModelError
	}
	return nil
}

func (c InventoryService) UpdateInventoryTopic(ctx context.Context, InventoryTopicUpdateRequest *request_dto.InventoryTopicUpdateRequest, TopicId string) *dto.ErrorResponseDto {
	methodName := "UpdateInventoryTopic"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto
	log.Info("Inside " + methodName)

	InventoryTopicUpdateModel, BindingError := utils.TypeConverter[models.InventoryTopicUpdateRequest](InventoryTopicUpdateRequest)

	if BindingError != nil {
		log.Info("Error while binding Remove Inventory Subject Model", BindingError)
		domainErr.SetError(status_code.IMS400)
		return &domainErr
	}

	UpdateInventoryTopicErr := c.InventoryRepository.UpdateInventoryTopic(ctx, InventoryTopicUpdateModel, TopicId)

	if UpdateInventoryTopicErr != nil {
		log.Info("Error while Updating Topic", UpdateInventoryTopicErr)
		return UpdateInventoryTopicErr
	}
	return nil
}

func (c InventoryService) ActivateResourceById(ctx context.Context, InventoryName string, Id string) *dto.ErrorResponseDto {
	methodName := "UpdateInventoryTopic"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)

	UpdateInventoryTopicErr := c.InventoryRepository.ActivateResourceById(ctx, InventoryName, Id)

	if UpdateInventoryTopicErr != nil {
		log.Info("Error while Updating Topic", UpdateInventoryTopicErr)
		return UpdateInventoryTopicErr
	}
	return nil
}

func (c InventoryService) CreateResource(ctx context.Context, InventoryResourceCreate request_dto.InventoryResourceCreate, topicId string, contentType string, fileExtension string) (*string, *dto.ErrorResponseDto) {
	methodName := "CreateResource"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto
	log.Info("Inside " + methodName)
	var resourceDto *Resource

	fileUrl := "https://" + viper.GetString("S3BucketName") + ".s3." + viper.GetString("aws.region") + ".amazonaws.com/kod/" + InventoryResourceCreate.ServiceName + "/" + InventoryResourceCreate.FlowType + "/" + topicId + "/" + InventoryResourceCreate.TopicName + "_" + topicId + "/" + InventoryResourceCreate.FileType + "/" + InventoryResourceCreate.FileType + "_" + InventoryResourceCreate.FileRequestId + fileExtension
	s3err := c.S3.UploadFile(InventoryResourceCreate.ServiceName, InventoryResourceCreate.FlowType, InventoryResourceCreate.TopicName, topicId, InventoryResourceCreate.Resource, InventoryResourceCreate.FileRequestId, contentType, fileExtension, InventoryResourceCreate.FileType)
	if s3err != nil {
		log.Error("Inside "+methodName+" error: ", s3err.Error(), " occurred while trying to upload resource to s3 for topic: ", topicId)
		domainErr.SetError(status_code.IMS400)
		return nil, &domainErr
	}

	ResourceModelObject, typeConvertError := utils.TypeConverter[ResourceModel](resourceDto)
	if typeConvertError != nil {
		log.Info("Error while binding Remove Inventory Subject Model", typeConvertError)
		domainErr.SetError(status_code.IMS400)
		return nil, &domainErr
	}
	ResourceModelObject.TopicId = topicId
	ResourceModelObject.Path = fileUrl
	ResourceModelObject.FlowType = InventoryResourceCreate.FlowType
	ResourceModelObject.ServiceName = InventoryResourceCreate.ServiceName
	ResourceModelObject.TopicName = InventoryResourceCreate.TopicName
	ResourceModelObject.Id = uuid.NewString()

	ResourceInterface, BindErr := utils.TypeConverter[interface{}](ResourceModelObject)
	if BindErr != nil {
		log.Info("Error while binding Remove Inventory Subject Model", BindErr)
		domainErr.SetError(status_code.IMS400)
		return nil, &domainErr
	}

	UpdateInventoryTopicErr := c.InventoryRepository.CreateNewInventoryGivenInventoryName(ctx, ResourceInterface, "resources")
	if UpdateInventoryTopicErr != nil {
		log.Info("Error while Updating Topic", UpdateInventoryTopicErr)
		return nil, UpdateInventoryTopicErr
	}
	log.Info("Update Err", UpdateInventoryTopicErr)

	return &fileUrl, nil
}

func KeyExists(list []request_dto.InventoryIdentifier, keyToFind string) bool {
	for _, item := range list {
		if item.Key == keyToFind {
			return true
		}
	}
	return false
}

func UniqueKeyExists(list []request_dto.InventoryIdentifier, keyToFind string) bool {
	for _, item := range list {
		if item.Key == keyToFind && item.IsUnique == true {
			return true
		}
	}
	return false
}

func (c InventoryService) GetInventoryFilter(ctx context.Context, InventoryName string, FilterName string, filters map[string][]string) ([]interface{}, *dto.ErrorResponseDto) {
	methodName := "FilterName"
	log := logger.GetLogger()
	var domainErr dto.ErrorResponseDto
	log.Info("Inside " + methodName)

	inventoryConfiguration, errDto := c.InventoryConfigurationService.GetInventoryConfiguration(ctx, InventoryName)
	if errDto != nil {
		if errDto.StatusCode == status_code.IMS204 {
			log.Info("Inside "+methodName+" inventory item deleted :", InventoryName)
			return nil, errDto
		}
		log.Info("Inside "+methodName+" unable to fetch inventory item for InventoryName :", InventoryName)
		return nil, errDto
	}

	//Checking if filter attribute exist in identifier list
	// log.Info("InventoryIdentifiers", inventoryConfiguration.InventoryIdentifiers, FilterName)
	if !KeyExists(inventoryConfiguration.InventoryIdentifiers, FilterName) {
		domainErr.SetError(status_code.IMS113)
		return nil, &domainErr
	}

	for key, _ := range filters {
		if !KeyExists(inventoryConfiguration.InventoryIdentifiers, key) {
			log.Error("Inside "+methodName+" filter attribute not present in inventoryIdentifiers :", InventoryName, " and filterMap : ", filters)
			domainErr.SetError(status_code.IMS113)
			return nil, &domainErr
		}
	}

	filtersBson := bson.M{}
	for key, value := range filters {
		filtersBson[key] = bson.M{"$in": value}
	}

	item, adapterError := c.InventoryRepository.GetInventoryFilter(ctx, InventoryName, FilterName, filtersBson)
	if adapterError != nil {
		log.Error("Inside "+methodName+" error while fetching item for :", InventoryName, " and filterMap : ", FilterName)
		return nil, adapterError
	}
	return item, nil

}

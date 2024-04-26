package impl

import (
	"context"
	"inventory-system/common/pkg/constants"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/inventory-service/internal/adapters/db"
	"inventory-system/inventory-service/internal/adapters/models"
	commonDto "inventory-system/inventory-service/internal/common/dto"
	"inventory-system/inventory-service/internal/common/status_code"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryRepository struct {
}

func NewInventoryRepository() *InventoryRepository {
	repo := &InventoryRepository{}
	return repo
}
func (c InventoryRepository) CreateNewInventoryGivenInventoryName(ctx context.Context, item interface{}, inventoryName string) *dto.ErrorResponseDto {
	methodName := "CreateNewInventoryGivenInventoryName"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto

	collectionName := constants.InventoryCollectionNamePrefix + inventoryName
	//collectionList, _ := db.GetDb().ListCollectionNames(ctx, bson.D{{}})
	//if !utils.Contains(collectionList, collectionName) {
	//	log.Error("Inside ", methodName, "collection name not found for ", collectionName)
	//	adapterErr.SetError(status_code.IMS404)
	//	return &adapterErr
	//}
	_, err := db.GetDb().Collection(collectionName).InsertOne(ctx, item)
	if err != nil {
		log.Error(err.Error())
		errorData := err.(mongo.ServerError)
		if errorData.HasErrorCode(constants.MongoDuplicateEntryErrorCode) {
			log.Error("Inside "+methodName+" error duplicate entry occurred when trying to create item ", item)
			adapterErr.SetError(status_code.IMS108)
			return &adapterErr
		}
		if strings.Contains(err.Error(), "Document failed validation") {
			log.Error("Inside ", methodName, "error : ", err.Error(), " occurred while creating item: ", item)
			adapterErr.SetError(status_code.IMS109)
			return &adapterErr
		}

		log.Error("Inside "+methodName+" error: ", err.Error(), " while creating the item: ", item)
		adapterErr.SetError(status_code.IMS101)
		return &adapterErr
	}

	return nil
}

func (c InventoryRepository) FetchInventory(ctx context.Context, inventoryName string, uniqueKey string, uniqueValue string) (bson.M, *dto.ErrorResponseDto) {
	methodName := "FetchInventory"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto
	collectionName := constants.InventoryCollectionNamePrefix + inventoryName

	filter := bson.M{uniqueKey: uniqueValue, "is_deleted": false}
	var item bson.M
	err := db.GetDb().Collection(collectionName).FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 0})).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("Inside "+methodName+" no documents exists in inventory for filter : ", filter)
			adapterErr.SetError(status_code.IMS404)
			return nil, &adapterErr
		}
		adapterErr.SetError(status_code.IMS110)
		return nil, &adapterErr

	}
	return item, nil
}

func (c InventoryRepository) RemoveItemFromInventory(ctx context.Context, RemoveItemModel *models.RemoveInventoryItem, InventoryName string) *dto.ErrorResponseDto {
	methodName := "RemoveItemFromInventory"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	var adapterErr dto.ErrorResponseDto
	log.Info("Request Item Id", RemoveItemModel.ItemId)

	log.Info("Inventory Name", InventoryName)
	collectionName := constants.InventoryCollectionNamePrefix + InventoryName
	log.Info("Collection Name", collectionName)
	isDeleted, Err := db.GetDb().Collection(collectionName).UpdateOne(ctx, bson.M{"id": RemoveItemModel.ItemId}, bson.M{"$set": bson.M{"is_deleted": true}})
	if Err != nil {
		log.Info("Error while removing the item from the Inventory with Inventory "+InventoryName, "And error is ", Err)
		adapterErr.SetError(status_code.IMS500)
		return &adapterErr
	}
	log.Info("Delete Feedback", isDeleted)
	return nil
}
func (c InventoryRepository) RemoveSubjectTopicsByLessonNameAndSubjectId(ctx context.Context, SubjectRemoveModel *models.RemoveSubjectRequestModel, Type string) *dto.ErrorResponseDto {
	methodName := "RemoveSubjectTopicsByLessonNameAndSubjectId"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	log.Info("Request Item Id", SubjectRemoveModel.SubjectId)
	log.Info("Type ", Type)

	if strings.ToLower(Type) == "partial" {
		removeErr := RemoveSubjectPartially(ctx, SubjectRemoveModel)
		return removeErr
	} else if strings.ToLower(Type) == "full" {
		FullRemoveErr := RemoveSubjectFull(ctx, SubjectRemoveModel)
		return FullRemoveErr
	}

	return nil
}

func RemoveSubjectPartially(ctx context.Context, SubjectRemoveModel *models.RemoveSubjectRequestModel) *dto.ErrorResponseDto {
	methodName := "RemoveSubjectPartially"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	var adapterErr dto.ErrorResponseDto
	collectionName := constants.InventoryCollectionNamePrefix + "topics"

	isDeleted, Err := db.GetDb().Collection(collectionName).DeleteMany(ctx, bson.M{"lesson_name": SubjectRemoveModel.LessonName, "subject_id": SubjectRemoveModel.SubjectId})
	if Err != nil {
		log.Info("Error while removing the item from the Inventory with Inventory And error is ", Err)
		adapterErr.SetError(status_code.IMS500)
		return &adapterErr
	}

	if isDeleted.DeletedCount != 1 {
		adapterErr.SetError(status_code.IMS118)
		return &adapterErr
	}

	return nil
}

func RemoveSubjectFull(ctx context.Context, SubjectRemoveModel *models.RemoveSubjectRequestModel) *dto.ErrorResponseDto {
	methodName := "RemoveSubjectFull"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	var adapterErr dto.ErrorResponseDto
	collectionName := constants.InventoryCollectionNamePrefix + "topics"
	isSubjectRemoved, DbErr := db.GetDb().Collection("Inventory-subjects").DeleteOne(ctx, bson.M{"id": SubjectRemoveModel.SubjectId})
	if DbErr != nil {
		log.Info("Error while removing Subject by its id", SubjectRemoveModel.SubjectId)
		adapterErr.SetError(status_code.IMS500)
		return &adapterErr
	}

	if isSubjectRemoved.DeletedCount != 1 {
		adapterErr.SetError(status_code.IMS116)
		return &adapterErr
	}
	if SubjectRemoveModel.LessonName == "" {
		isDeleted, Err := db.GetDb().Collection(collectionName).DeleteMany(ctx, bson.M{"subject_id": SubjectRemoveModel.SubjectId})
		if Err != nil {
			log.Info("Error while removing the item from the Inventory with Inventory And error is ", Err)
			adapterErr.SetError(status_code.IMS500)
			return &adapterErr
		}

		if isDeleted.DeletedCount != 1 {
			adapterErr.SetError(status_code.IMS116)
			return &adapterErr
		}
	}

	return nil
}

func (c InventoryRepository) UpdateInventoryTopic(ctx context.Context, InventoryTopicUpdateModel *models.InventoryTopicUpdateRequest, TopicId string) *dto.ErrorResponseDto {
	methodName := "UpdateInventoryTopic"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	var adapterErr dto.ErrorResponseDto
	log.Info("Topic Update request", InventoryTopicUpdateModel)
	collectionName := constants.InventoryCollectionNamePrefix + "topics"
	log.Info("Collection Name", collectionName)

	log.Info("Topic Id", TopicId)
	Update, DbErr := db.GetDb().Collection(collectionName).UpdateOne(ctx, bson.M{"id": TopicId}, bson.M{"$set": bson.M{"resources": InventoryTopicUpdateModel.Resources, "topic_name": InventoryTopicUpdateModel.TopicName, "description": InventoryTopicUpdateModel.Description, "assessments": InventoryTopicUpdateModel.Assessments}})
	if DbErr != nil {
		log.Info("Error while Updating Topic with topic id", TopicId)
		adapterErr.SetError(status_code.IMS500)
		return &adapterErr
	}

	log.Info("Update", Update)

	if Update.ModifiedCount != 1 && Update.MatchedCount != 1 {

		adapterErr.SetError(status_code.IMS117)
		return &adapterErr
	}

	log.Info("Update Feedback", Update)
	return nil
}

func (c InventoryRepository) UpdateInventory(Id string, InventoryName string, UpdateRequest *interface{}) *dto.ErrorResponseDto {
	methodName := "UpdateInventoryTopic"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	var adapterErr dto.ErrorResponseDto
	collectionName := constants.InventoryCollectionNamePrefix + InventoryName
	log.Info("Collection Name", collectionName)
	updateRes, dbErr := db.GetDb().Collection(collectionName).UpdateOne(context.Background(), bson.M{"id": Id}, bson.M{"$set": UpdateRequest})
	if dbErr != nil {
		log.Info("There is an error while updating the inventory", dbErr)
		adapterErr.SetError(status_code.IMS101)
		return &adapterErr
	}

	log.Info(updateRes)
	return nil
}

func (c InventoryRepository) ActivateResourceById(ctx context.Context, InventoryName string, Id string) *dto.ErrorResponseDto {
	methodName := "ActivateResourceById"
	log := logger.GetLogger()
	log.Info("Inside " + methodName)
	var adapterErr dto.ErrorResponseDto
	log.Info("Inventory Name", InventoryName)
	collectionName := constants.InventoryCollectionNamePrefix + InventoryName
	log.Info("Collection Name", collectionName)

	log.Info("Id", Id)
	Update, DbErr := db.GetDb().Collection(collectionName).UpdateOne(ctx, bson.M{"id": Id}, bson.M{"$set": bson.M{"is_deleted": false}})
	if DbErr != nil {
		log.Info("Error while Updating with id", Id)
		adapterErr.SetError(status_code.IMS600)
		return &adapterErr
	}

	log.Info("Update", Update)

	if Update.ModifiedCount != 1 && Update.MatchedCount != 1 {

		adapterErr.SetError(status_code.IMS119)
		return &adapterErr
	}

	log.Info("Update Feedback", Update)
	return nil
}

func (c InventoryRepository) FetchInventoryList(ctx context.Context, from string, to string, inventoryName string, filterMap map[string][]string, pagination commonDto.Pagination) ([]bson.M, *commonDto.PaginationResponse, *dto.ErrorResponseDto) {
	methodName := "FetchInventoryList"
	log := logger.GetLogger()
	var adapterErr dto.ErrorResponseDto
	collectionName := constants.InventoryCollectionNamePrefix + inventoryName

	query := bson.M{"is_deleted": false}

	if from != "" && to != "" {
		log.Info("FROM", from)
		log.Info("TO", to)

		query = bson.M{"created_at": bson.M{"$gte": from, "$lte": to}}

	}

	log.Info("Filter Map", filterMap)
	for key, values := range filterMap {
		log.Info("key")
		if len(values) == 1 {
			query[key] = values[0]
		} else {
			query[key] = bson.M{"$in": values}
		}
	}

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}}).SetProjection(bson.M{"_id": 0})

	if pagination.Pagination {
		opts.SetSkip(int64(pagination.PageSize) * int64(pagination.PageNumber)).SetLimit(int64(pagination.PageSize))
	}

	log.Info("Query", query)
	var itemList []bson.M
	cur, err := db.GetDb().Collection(collectionName).Find(ctx, query, opts)
	if err != nil {
		log.Info("Inside "+methodName+" error inventory for filter : ", query)
		adapterErr.SetError(status_code.IMS110)
		return nil, nil, &adapterErr
	}

	err = cur.All(ctx, &itemList)
	if err != nil {
		log.Error("Inside " + methodName + " error while decoding all inventory items")
		adapterErr.SetError(status_code.IMS306)
		return nil, nil, &adapterErr

	}
	if itemList == nil {
		log.Error("Inside " + methodName + " no document exists")
		adapterErr.SetError(status_code.IMS404)
		return nil, nil, &adapterErr
	}

	count, err2 := db.GetDb().Collection(collectionName).CountDocuments(ctx, query)
	if err2 != nil {
		log.Error("Inside " + methodName + " error while counting all inventory items")
		adapterErr.SetError(status_code.IMS306)
		return nil, nil, &adapterErr
	}

	paginationResp := &commonDto.PaginationResponse{}

	if pagination.Pagination {
		paginationResp = &commonDto.PaginationResponse{
			Count:      count,
			PageNumber: pagination.PageNumber,
		}
	}
	if pagination.Pagination {
		return itemList, paginationResp, nil
	}
	return itemList, nil, nil
}

func (c InventoryRepository) GetInventoryFilter(ctx context.Context, InventoryName string, FilterName string, filters bson.M) ([]interface{}, *dto.ErrorResponseDto) {
	methodName := "GetInventoryFilter"
	log := logger.GetLogger()

	var adapterErr dto.ErrorResponseDto
	collectionName := constants.InventoryCollectionNamePrefix + InventoryName

	collection := db.GetDb().Collection(collectionName)
	filterList, err := collection.Distinct(ctx, FilterName, filters)
	if err != nil {
		log.Info(err.Error())
		log.Info("Inside "+methodName+" error inventory for filter : ", FilterName)
		adapterErr.SetError(status_code.IMS110)
		return nil, &adapterErr
	}

	return filterList, nil
}

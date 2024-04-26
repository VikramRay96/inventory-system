package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"inventory-system/common/pkg/dto"
	"inventory-system/common/pkg/logger"
	"inventory-system/inventory-service/internal/common/dto/request_dto"
	"inventory-system/inventory-service/internal/common/status_code"
	"inventory-system/inventory-service/internal/domain/service"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	InventoryService service.IInventoryService
}

func NewInventoryController(inventoryService service.IInventoryService) *InventoryController {
	controller := &InventoryController{
		InventoryService: inventoryService,
	}
	return controller
}

// AddNewInventory  godoc
// @Summary Add item to inventory
// @Description Add item to inventory
// @Tags Inventory
// @Produce  json
// @Success 200 {object} dto.ResponseDto
// @Param requestBody body bson.M true "Add item to inventory request"
// @Param inventoryName path string true "Inventory Key"
// @Router /inventory-service/api/v1/inventory/configurations/{inventoryName} [POST]
// AddNewInventory : This function will an item from inventory
func (cc InventoryController) AddNewInventory() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "AddNewInventory"
		log := logger.GetLogger()
		ctx := context.Background()
		inventoryName := c.Param("inventoryName")
		var requestBody interface{}
		var portErr dto.ErrorResponseDto

		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			log.Info(err)
			portErr.SetError(status_code.IMS400)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: portErr.StatusCode,
				Message:    portErr.Message,
			})
			return
		}
		errorDto := cc.InventoryService.CreateNewInventory(ctx, requestBody, inventoryName)
		if errorDto != nil {
			log.Error("Inside "+methodName+" error while adding inventory: ", requestBody)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errorDto.StatusCode,
				Message:    errorDto.Message,
			})
			return
		}
		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})
	}
	return fn
}

// GetInventory  godoc
// @Summary Get an item from inventory
// @Description Get Inventory
// @Tags Inventory
// @Produce  json
// @Success 200 {object} dto.ResponseDto
// @Param inventoryName path string true "Inventory Key"
// @Param {inventory_identifier_key} query string true "inventory Identifier Value"
// @Router /inventory-service/api/v1/{inventoryName}/inventory [GET]
// GetInventory : This function will fetch an item from inventory
func (cc InventoryController) GetInventory() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "GetInventory"
		log := logger.GetLogger()
		ctx := context.Background()
		filterMap := c.Request.URL.Query()

		inventoryConfigurationName := c.Param("inventoryName")

		inventory, errDto := cc.InventoryService.GetInventory(ctx, inventoryConfigurationName, filterMap)
		if errDto != nil {
			log.Info("Inside "+methodName+" unable to fetch inventory for inventoryConfigurationName :", inventoryConfigurationName, " filterMap: ", filterMap)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}
		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			Data:       inventory,
		})
	}
	return fn
}

// GetInventoryV2  godoc
// @Summary Get an item from inventory
// @Description Get Inventory
// @Tags Inventory
// @Produce  json
// @Success 200 {object} dto.ResponseDto
// @Param inventoryName path string true "Inventory Key"
// @Param {inventory_identifier_key} query string true "inventory Identifier Value"
// @Router /inventory-service/api/v1/{inventoryName}/inventory [GET]
// GetInventory : This function will fetch an item from inventory
func (cc InventoryController) GetInventoryV2() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		methodName := "GetInventory"
		log := logger.GetLogger()
		// ctx := context.Background()
		from := c.Query("from")
		to := c.Query("to")

		inventoryConfigurationName := c.Param("inventoryName")
		var filterMap map[string][]string
		var portErr dto.ErrorResponseDto
		err := c.ShouldBindJSON(&filterMap)
		if err != nil {

			log.Info("Filter Err")
			portErr.SetError(status_code.IMS400)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: portErr.StatusCode,
				Message:    portErr.Message,
			})
			return
		}
		inventory, pagination, errDto := cc.InventoryService.GetInventoryV2(c, from, to, inventoryConfigurationName, filterMap)
		if errDto != nil {
			log.Info("Inside "+methodName+" unable to fetch inventory for inventoryConfigurationName :", inventoryConfigurationName, " filterMap: ", filterMap)
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
				Data:       []bson.M{},
			})
			return
		}
		data := []bson.M{}
		if len(inventory) > 0 {
			data = inventory
		}

		log.Info("Pagination", pagination)

		if pagination != nil {
			c.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
				Message:    dto.GetStatusDetails(status_code.IMS200).Message,
				Data: bson.M{
					"count": pagination.Count,
					"page":  pagination.PageNumber,
					"items": data,
				},
			})
			return
		}

		c.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			Data:       data,
		})
	}
	return fn
}

func (cc InventoryController) RemoveItemFromInventory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "RemoveItemFromInventory"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)
		var errorDto *dto.ErrorResponseDto
		var RemoveItemRequest *request_dto.RemoveInventoryItem
		InventoryName := ctx.Param("inventoryName")

		err := ctx.ShouldBindJSON(&RemoveItemRequest)
		if err != nil {
			errorDto.SetError(status_code.IMS400)
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errorDto.StatusCode,
				Message:    errorDto.Message,
			})
			return
		}

		log.Info("RemoveItemRequest", RemoveItemRequest)

		errDto := cc.InventoryService.RemoveItemFromInventory(ctx, RemoveItemRequest, InventoryName)
		if errDto != nil {
			log.Info("Inside " + methodName + " unable to fetch inventory for inventoryConfigurationName :")
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})

	}
}

func (cc InventoryController) UpdateInventory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "UpdateInventory"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)
		var UpdateRequest *interface{}
		InventoryName := ctx.Param("inventoryName")
		Id := ctx.Param("id")

		bindErr := ctx.ShouldBind(&UpdateRequest)
		if bindErr != nil {
			log.Info("Invalid Request Body", bindErr)
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: dto.GetStatusDetails(status_code.IMS400).StatusCode,
				Message:    dto.GetStatusDetails(status_code.IMS400).Message,
			})
			return
		}

		errorDto := cc.InventoryService.UpdateInventory(Id, InventoryName, UpdateRequest)
		if errorDto != nil {
			log.Info("There is an issue while updating Inventory", errorDto)
			ctx.JSON(http.StatusOK, errorDto)
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Data:       "Inventory Updated Successfully",
		})
	}
}

func (cc InventoryController) GetInventoryFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "UpdateInventory"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)
		InventoryName := ctx.Param("inventoryName")
		FilterName := ctx.Param("filterName")
		filters := ctx.Request.URL.Query()

		data, errorDto := cc.InventoryService.GetInventoryFilter(ctx, InventoryName, FilterName, filters)
		if errorDto != nil {
			log.Info("There is an issue while get Inventory", errorDto)
			ctx.JSON(http.StatusOK, errorDto)
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Data:       data,
		})
	}
}

func (cc InventoryController) RemoveSubjectTopicsByLessonNameAndSubjectId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "RemoveSubjectTopicsByLessonNameAndSubjectId"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)
		var errorDto *dto.ErrorResponseDto
		var RemoveSubjectRequest *request_dto.RemoveSubjectRequest
		Type := ctx.Query("type")

		log.Info(strings.ToLower(Type) != "full")

		if strings.ToLower(Type) != "partial" && strings.ToLower(Type) != "full" {
			log.Info("Inside IF", Type)
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: dto.GetStatusDetails(status_code.IMS400).StatusCode,
				Message:    dto.GetStatusDetails(status_code.IMS400).Message,
			})
			return
		}

		err := ctx.ShouldBindJSON(&RemoveSubjectRequest)
		if err != nil {
			errorDto.SetError(status_code.IMS400)
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errorDto.StatusCode,
				Message:    errorDto.Message,
			})
			return
		}

		errDto := cc.InventoryService.RemoveSubjectTopicsByLessonNameAndSubjectId(ctx, RemoveSubjectRequest, Type)
		if errDto != nil {
			log.Info("Inside " + methodName + " unable to fetch inventory for inventoryConfigurationName :")
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})
	}
}

func (cc InventoryController) UpdateInventoryTopic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "UpdateInventoryTopic"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)
		var errorDto *dto.ErrorResponseDto
		var InventoryTopicUpdate *request_dto.InventoryTopicUpdateRequest
		TopicId := ctx.Param("topic_id")

		err := ctx.ShouldBindJSON(&InventoryTopicUpdate)
		if err != nil {
			errorDto.SetError(status_code.IMS400)
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errorDto.StatusCode,
				Message:    errorDto.Message,
			})
			return
		}

		errDto := cc.InventoryService.UpdateInventoryTopic(ctx, InventoryTopicUpdate, TopicId)
		if errDto != nil {
			log.Info("Inside " + methodName + " unable to update Inventory Topic")
			ctx.JSON(http.StatusOK, dto.ResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})
	}
}

func (cc InventoryController) ActivateResourceById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "ActivateResourceById"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)
		InventoryName := ctx.Param("inventoryName")
		Id := ctx.Query("id")

		errDto := cc.InventoryService.ActivateResourceById(ctx, InventoryName, Id)
		if errDto != nil {
			log.Info("Inside " + methodName + " unable to update Inventory Topic")
			ctx.JSON(http.StatusOK, dto.ErrorResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
		})
	}
}

func (cc InventoryController) CreateResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		methodName := "CreateResource"
		log := logger.GetLogger()
		log.Info("Inside " + methodName)

		form, err := ctx.MultipartForm()
		if err != nil {
			log.Info("Inside Error", err)
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		serviceName := ctx.PostForm("service_name")
		flowType := ctx.PostForm("flow_type")
		topicId := ctx.PostForm("topic_id")
		fileType := ctx.PostForm("file_type")
		createdBy := ctx.PostForm("created_by")
		topicName := ctx.PostForm("topic_name")
		file_request_id := ctx.PostForm("file_request_id")

		// Handle file
		files := form.File["resource"]
		if len(files) != 1 {
			ctx.JSON(http.StatusOK, gin.H{"error": "Bad Request"})
			return
		}

		videoFile, err := files[0].Open()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		defer videoFile.Close()

		contentType := files[0].Header.Get("Content-Type")

		fileExtension := filepath.Ext(files[0].Filename)

		log.Info("File Type", contentType)
		log.Info("File Extension ", fileExtension)

		videoBytes, err := io.ReadAll(videoFile)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		InventoryResourceCreate := request_dto.InventoryResourceCreate{
			FlowType:      flowType,
			ServiceName:   serviceName,
			TopicId:       topicId,
			Resource:      videoBytes,
			FileType:      fileType,
			CreatedBy:     createdBy,
			TopicName:     topicName,
			FileRequestId: file_request_id,
		}

		log.Info("InventoryResourceCreate", InventoryResourceCreate.CreatedBy)
		response, errDto := cc.InventoryService.CreateResource(ctx, InventoryResourceCreate, topicId, contentType, fileExtension)
		if errDto != nil {
			log.Info("Inside " + methodName + " unable to update Inventory Topic")
			ctx.JSON(http.StatusOK, dto.ErrorResponseDto{
				StatusCode: errDto.StatusCode,
				Message:    errDto.Message,
			})
			return
		}

		ctx.JSON(http.StatusOK, dto.ResponseDto{
			StatusCode: dto.GetStatusDetails(status_code.IMS200).StatusCode,
			Message:    dto.GetStatusDetails(status_code.IMS200).Message,
			Data:       response,
		})
	}
}

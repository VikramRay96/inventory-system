package router

import (
	commonConstants "inventory-system/common/pkg/constants"
	"inventory-system/common/pkg/dto"
	"inventory-system/inventory-service/internal/common/status_code"
	portConstants "inventory-system/inventory-service/internal/ports/constants"
	"inventory-system/inventory-service/internal/ports/docs"
	"inventory-system/inventory-service/internal/ports/factory"
	"strings"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"bitbucket.org/kodnest/go-common-libraries/correlation"
	"bitbucket.org/kodnest/go-common-libraries/logger"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/spf13/viper"
)

func NewRouter() *gin.Engine {

	if strings.EqualFold(viper.GetString(commonConstants.ENVIRONMENT), commonConstants.ENV_PRODUCTION) {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(func(c *gin.Context) {
	defer func() {
	if r := recover(); r != nil {
	c.JSON(500, dto.ErrorResponseDto{
	StatusCode: dto.GetStatusDetails(status_code.IMS500).StatusCode,
	Message:    dto.GetStatusDetails(status_code.IMS500).Message,
	})
	}
		}()

	c.Next()
	})
	//
	//allowedOrigins := viper.GetString(commonConstants.CORS_ALLOWED_ORIGINS)
	//allowedOriginsList := strings.Split(allowedOrigins, ",")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	})

	controllerFacade := factory.GetControllerFacade()
	router.Use(corsHandler)
	router.GET(portConstants.HEALTH_CHECK_PATH, controllerFacade.HealthController.Status())

	router.Use(uuidInjectionMiddleware())

	//swag init --parseDependency --parseInternal -g internal/ports/router/router.go
	if !strings.EqualFold(commonConstants.ENVIRONMENT, commonConstants.ENV_PRODUCTION) {

		docs.SwaggerInfo.Title = "Inventory Service"
		docs.SwaggerInfo.Version = "1.8.7"
		docs.SwaggerInfo.Description = "Swagger UI for configuration service"

		swaggerUrl := viper.GetString(commonConstants.BASE_URL) + commonConstants.SERVICE_NAME + commonConstants.SwaggerDocPath
		url := ginSwagger.URL(swaggerUrl) // The url pointing to API definition
		router.GET(commonConstants.BASE_PATH+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}

	//Implementation Of Swagger Needs To Be Done
	configurationService := router.Group(commonConstants.BASE_PATH)
	{
		api := configurationService.Group(portConstants.API)
		{
			v1 := api.Group(portConstants.VERSION_V1)
			{
				//Base Configuration
				v1.GET("/inventory/configurations/:inventoryName", controllerFacade.InventoryConfigurationController.GetConfiguration())
				v1.GET("/inventory/configurations", controllerFacade.InventoryConfigurationController.GetAllConfiguration())
				v1.DELETE("/inventory/configurations/:inventoryName", controllerFacade.InventoryConfigurationController.DeleteConfiguration())
				v1.POST("/inventory/configurations", controllerFacade.InventoryConfigurationController.CreateNewConfiguration())
				//Inventory Controller
				v1.PATCH("/inventory/:inventoryName", controllerFacade.InventoryController.ActivateResourceById())
				v1.GET("/inventory/:inventoryName", controllerFacade.InventoryController.GetInventory())
				v1.POST("/inventory/:inventoryName", controllerFacade.InventoryController.AddNewInventory())
			}
			v2 := api.Group(portConstants.VERSION_V2)
			{
				v2.PATCH("/inventory/update/:inventoryName/:id", controllerFacade.InventoryController.UpdateInventory())
				v2.GET("/inventory/filter/:inventoryName/:filterName", controllerFacade.InventoryController.GetInventoryFilter())
				v2.POST("/inventory/:inventoryName", controllerFacade.InventoryController.GetInventoryV2())
				v2.DELETE("/inventory/:inventoryName", controllerFacade.InventoryController.RemoveItemFromInventory())
				v2.DELETE("/inventory/subject/remove", controllerFacade.InventoryController.RemoveSubjectTopicsByLessonNameAndSubjectId())

				v2.PATCH("/inventory/topics/:topic_id/update", controllerFacade.InventoryController.UpdateInventoryTopic())
				v2.POST("/inventory/resource/create", controllerFacade.InventoryController.CreateResource())
			}

		}
	}

	return router
}

func uuidInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.New(logger.Info)
		correlationId, err := correlation.FromContext(c)
		if err != nil {
			log.Info("CorrelationId not found in gin context, creating new correlationId")
		}
		if len(correlationId) == 0 {
			correlationId := correlation.NewId()
			c.Request.Header.Set(commonConstants.CorrelationId, correlationId)
		}
		c.Writer.Header().Set(commonConstants.CorrelationId, correlationId)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

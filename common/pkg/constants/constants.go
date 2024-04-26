package constants

const (
	SERVICE_NAME = "inventory-service"
)

const (
	COLON                      = ":"
	ENV_PRODUCTION             = "production"
	ENV_STAGING                = "staging"
	ENV_DEV                    = "dev"
	BOOT_CUR_ENV               = "BOOT_CUR_ENV"
	UNDERSCORE_SEPARATOR       = "_"
	ENV_IS_VERSION             = "version"
	VERSION_FILE_NAME          = "VERSION"
	ERROR_KEY_ATTRIBUTE_NAME   = "error"
	HYPEN_SEPARATOR            = "-"
	ERROR_ATTRIBUTE_NAME       = "ERROR"
	CORRELATION_ATTRIBUTE_NAME = "correlationId"
	SERVICE_NAME_KEY           = "serviceName"
	CorrelationId              = "X-Correlation-ID"
	BASE_URL                   = "BaseUrl"
	SwaggerDocPath             = "/swagger/doc.json"
)

// Config file path
const (
	CONFIG_FILE_NAME                   = "/config/config.json"
	SERVER_PRODUCTION_CONFIG_FILE_NAME = "/configuration/production/config.json"
	SERVER_STAGING_CONFIG_FILE_NAME    = "/configuration/staging/config.json"
)

// ViperKeys
const (
	CORS_ALLOWED_ORIGINS = "allowedOrigins"
	ENVIRONMENT          = "Environment"
	AWS_KEY              = "AWS_KEY"
	AWS_SECRET           = "AWS_SECRET"

	MAX_OPTMISTIC_LOCKING_RETRY_COUNT = "maxOptimisticLockingRetryCount"

	//rest config
	REST_EXECUTE_TIME_OUT_IN_SEC = "RestExecuteTimeoutInSeconds"

	//Retry config
	MESSAGE_PUBLISHER_RETRY_COUNT = "pls.sns.publisher.retryCount"

	BASE_PATH = "inventory-service"
)

// Redis
const (
	CACHE_TTL = "cacheTtl"
)

const (
	MAX_STARTUP_ATTEMPT = "pls.maxStartupAttempt"
)
const (
	LoggerLevelKey = "log.level"
)

const (
	InventoryConfigurationCollectionName = "InventoryConfiguration"
	MongoDuplicateEntryErrorCode         = 11000
	InventoryCollectionNamePrefix        = "Inventory-"
)

package config

import (
	commonConstants "inventory-system/common/pkg/constants"
)

var PropertiesMap = map[string]interface{}{
	commonConstants.MAX_OPTMISTIC_LOCKING_RETRY_COUNT: 3,
	commonConstants.MAX_STARTUP_ATTEMPT:               3,
}

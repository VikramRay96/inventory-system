package logger

import (
	"bitbucket.org/kodnest/go-common-libraries/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"inventory-system/common/pkg/constants"
)

// GetLogger : Return SugaredLogger instance from library according to log level set in config
// kept variadic fn as context is optional parameter
func GetLogger() *zap.SugaredLogger {
	logLevel := logger.LogLevel(viper.GetInt(constants.LoggerLevelKey))
	return logger.New(logLevel)
}

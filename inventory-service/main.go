package main

import (
	"bitbucket.org/kodnest/go-common-libraries/config"
	"bitbucket.org/kodnest/go-common-libraries/gracefulshutdown"
	"bitbucket.org/kodnest/go-common-libraries/logger"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	commonConstants "inventory-system/common/pkg/constants"
	"inventory-system/common/pkg/utils"
	configs "inventory-system/inventory-service/config"
	"inventory-system/inventory-service/internal/adapters/db"
	portConstants "inventory-system/inventory-service/internal/ports/constants"
	"inventory-system/inventory-service/pkg"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	log := logger.New(logger.Info)
	service := commonConstants.SERVICE_NAME
	environment := os.Getenv(commonConstants.BOOT_CUR_ENV)
	aws_key := os.Getenv("AWS_KEY")
	aws_secret := os.Getenv("AWS_SECRET")
	if environment == "" {
		environment = commonConstants.ENV_DEV
	}

	_, mainFilePath, _, _ := runtime.Caller(0)
	projectRootDir := filepath.Dir(mainFilePath)

	if ver, err := utils.NewServiceUtils().SetApplicationVersion(commonConstants.ENV_IS_VERSION, projectRootDir+"/"+commonConstants.VERSION_FILE_NAME); err != nil {
		log.Errorf("error setting version env versions  := %s , err := %v", ver, err)
	}
	version := os.Getenv(commonConstants.ENV_IS_VERSION)
	viper.Set(commonConstants.ENVIRONMENT, environment)
	viper.Set(commonConstants.AWS_KEY, aws_key)
	viper.Set(commonConstants.AWS_SECRET, aws_secret)
	log.Info("Version ", version)
	log.Info("Environment ", environment)
	log.Info("Service Name", service)
	configuration := config.New()

	if environment == commonConstants.ENV_DEV {
		configuration = configuration.FromFile(projectRootDir + commonConstants.CONFIG_FILE_NAME)
		if configuration.HasErrors() {
			log.Error(configuration.Errors)
			log.Error("error while fetching configurations from local config file Error: ", configuration.Errors)
			os.Exit(1)
		}
	} else if environment == commonConstants.ENV_PRODUCTION {
		log.Info("Taking environment config from ", "Production")
		configuration = configuration.FromFile(projectRootDir + commonConstants.SERVER_PRODUCTION_CONFIG_FILE_NAME)
		log.Info("CLOUD CONFIG", configuration)
		if configuration.HasErrors() {
			log.Error("error while fetching configurations from spring cloud Error: ", configuration.Errors)
			os.Exit(1)
		}
	} else if environment == commonConstants.ENV_STAGING {
		log.Info("Taking environment config from ", "staging")
		configuration = configuration.FromFile(projectRootDir + commonConstants.SERVER_STAGING_CONFIG_FILE_NAME)
		log.Info("CLOUD CONFIG", configuration)
		if configuration.HasErrors() {
			log.Error("error while fetching configurations from spring cloud Error: ", configuration.Errors)
			os.Exit(1)
		}
	}
	configuration.SetViper()
	log.Info(viper.AllSettings())
}

func main() {
	log := logger.New(logger.Info)
	exitChannel := make(chan struct{})
	gracefulShutDownManager := gracefulshutdown.NewManager(log, exitChannel)
	utils.NewServiceUtils().SetDefaultProperties(configs.PropertiesMap)
	db.Init()
	startRestApiService(gracefulShutDownManager, exitChannel)

}

// initializes all dependencies ,gracefull shutdown ,routers and starts a rest api service
func startRestApiService(gracefulShutDownManager *gracefulshutdown.Manager, exitChannel chan struct{}) {
	log := logger.New(logger.Info)
	flag.Usage = func() {
		fmt.Println("Usage: server -s {service_name} -e {environment}")
		os.Exit(1)
	}
	flag.Parse()
	/* Init gracefulshutdown */

	appRouter := pkg.RouterDriver()

	port := viper.GetString(portConstants.SERVER_PORT)
	httpServer := &http.Server{
		Addr:    port,
		Handler: appRouter,
	}
	/*
	   Starting the Http Server
	*/
	utils.SafeGoRoutine(commonConstants.SERVICE_NAME, func() {
		log.Info("Server Starting on Port : ", port)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal("listen: ", port)
		}
	}, viper.GetInt(commonConstants.MAX_STARTUP_ATTEMPT))

	// Init Shutdown Signals & Actions
	gracefulShutDownManager.Shutdown(httpServer)
	// Blocking until the shutdown to complete then inform the main goroutine.
	<-exitChannel
	log.Info("main goroutine shutdown completed gracefully.")
}

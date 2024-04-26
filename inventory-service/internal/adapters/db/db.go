package db

import (
	"fmt"
	"inventory-system/common/pkg/logger"
	"path/filepath"
	"runtime"
	"strings"

	db "bitbucket.org/kodnest/go-common-libraries/db/mongo"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var database *mongo.Database

type SecretPayload struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

// Init : Initializes the database migration
func Init() {
	log := logger.GetLogger()

	// sm := utils.NewSecretManager("ap-south-1", viper.GetString("secret_name"), viper.GetString("AWS_KEY"), viper.GetString("AWS_SECRET"))

	// secret, e := sm.GetSecrets()
	// if e != nil {
	// 	log.Info("Error while getting secret", e)
	// }

	// var dbConfig SecretPayload
	// err := json.Unmarshal([]byte(secret), &dbConfig)
	// if err != nil {
	// 	log.Fatalf("Error occurred during unmarshaling. Error: %s", err.Error())
	// }
	// dbUrl := fmt.Sprintf(viper.GetString("database.url"), dbConfig.UserName, dbConfig.Password, dbConfig.Host)
	dbUrl := fmt.Sprintf(viper.GetString("database.url"), viper.GetString("database.username"), viper.GetString("database.password"), viper.GetString("database.host"))
	log.Info("DB URL", dbUrl)
	viper.Set("mongo.uri", dbUrl)
	log.Info(viper.GetString("mongo.uri"))
	client, err := db.NewClient(log, viper.GetString("mongo.uri"))
	if err != nil {
		log.Fatal("Error ", err.Error(), " error connecting mongo db")
	}
	database = client.Database(viper.GetString("database.name"))
	driver, err := mongodb.WithInstance(client, &mongodb.Config{
		DatabaseName: viper.GetString("database.name"),
	})
	if err != nil {
		log.Fatal("Error ", err.Error(), " creating mongo driver instance")
	}
	_, mainFilePath, _, _ := runtime.Caller(0)
	projectRootDir := filepath.Dir(mainFilePath)
	migration, err := migrate.NewWithDatabaseInstance("file:"+projectRootDir+"/migrations", viper.GetString("database.name"), driver)
	if err != nil {
		log.Fatal("Error ", err.Error(), " while creating mongo migration instance")
	}
	currentActiveVersion, _, _ := migration.Version()
	log.Info("Current Active Version Of Migration :-", currentActiveVersion)
	err = migration.Up()
	if err != nil {
		log.Error("Error ", err.Error(), " while trying to run migration for version: ", currentActiveVersion)
		if strings.Contains(err.Error(), "Fix and force version") {
			err = migration.Force(int(currentActiveVersion))
			if err != nil {
				log.Fatal("Error ", err.Error(), " storing the current migration version: ", currentActiveVersion)
				return
			}

		}
	}

}

func GetDb() *mongo.Database {
	return database
}

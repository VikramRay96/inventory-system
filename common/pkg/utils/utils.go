package utils

import (
	"bitbucket.org/kodnest/go-common-libraries/logger"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
)

type Utils struct {
}

var (
	utilsService      *Utils //global variable for saving instance of Utils
	once              sync.Once
	valueReplaceRegex = regexp.MustCompile(`(?m)(\$\{(?P<placeholderKey>.*)\})`)
)

type IUtils interface {
	SetApplicationVersion(key, filename string) (string, error)
}

// NewServiceUtils - Creates new service utils object
func NewServiceUtils() *Utils {
	if utilsService == nil {
		once.Do(func() {
			utilsService = &Utils{}
		})
	}
	return utilsService
}

func (u *Utils) SetDefaultProperties(propertiesMap map[string]interface{}) {
	log := logger.New(logger.Warn)
	for key, value := range propertiesMap {
		if !viper.IsSet(key) {
			log.Warnf(fmt.Sprintf("no value set for property :  %s", key))
			viper.SetDefault(key, value)
		}
	}
}

func (u *Utils) SetApplicationVersion(key, filename string) (string, error) {
	var (
		version []byte
		err     error
	)
	version, err = ioutil.ReadFile(filename)
	os.Setenv(key, string(version))
	return string(version), err
}

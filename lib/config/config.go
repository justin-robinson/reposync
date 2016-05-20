package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"io"
	"log"
	"fmt"
)

const CONFIG_FILE_NAME = ".reposync.json"

type ReposyncConfig struct {
	Source struct {
			Repos map[string]bool
			Token string
			Type string
			Url string
	       }
}

// our config singleton
var instance *ReposyncConfig
var once sync.Once

var configExists = false

func ConfigExists() bool {
	return configExists
}

// reads and parses .reposync.json file inside of dir
func GetConfig(dir string) *ReposyncConfig {

	once.Do(func() {
		configFilePath := filepath.Join(dir, CONFIG_FILE_NAME)

		// open the config file
		configFile, configError := os.Open(configFilePath)

		// check errors to see if config file exists
		configExists = configError == nil

		// parse config if it exists
		if configExists {
			jsonDecoder := json.NewDecoder(configFile)

			// log error on json decode
			if err := jsonDecoder.Decode(&instance); err != io.EOF && err != nil {
				log.Fatal(fmt.Sprintf("Error decoding %v : %v", configFilePath, err.Error()))
			}
		} else {
			instance = &ReposyncConfig{}
		}
	})

	return instance
}

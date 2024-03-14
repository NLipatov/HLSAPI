package Configuration

import (
	"encoding/json"
	"os"
)

var _configurationPath = ""

type Configuration struct {
	Port              int    `json:"port"`
	StorageFolderPath string `json:"storageFolderPath"`
}

func Init(configurationPath string) {
	_configurationPath = configurationPath
}

func ReadConfiguration() Configuration {
	if len(_configurationPath) == 0 {
		panic("configuration manager was not initialized.")
	}
	configBytes, err := os.ReadFile(_configurationPath)
	if err != nil {
		panic(err)
	}

	config := Configuration{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}

package Configuration

import (
	"encoding/json"
	"hlsapi/src/Configuration/Models"
	"os"
)

var _configurationPath = ""

func Init(configurationPath string) {
	_configurationPath = configurationPath
}

func ReadConfiguration() ConfigurationModels.ConfigurationRoot {
	if len(_configurationPath) == 0 {
		panic("configuration manager was not initialized.")
	}
	configBytes, err := os.ReadFile(_configurationPath)
	if err != nil {
		panic(err)
	}

	config := ConfigurationModels.ConfigurationRoot{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}

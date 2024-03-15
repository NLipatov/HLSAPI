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

func ReadConfiguration() ConfigurationModels.Configuration {
	if len(_configurationPath) == 0 {
		panic("configuration manager was not initialized.")
	}
	configBytes, err := os.ReadFile(_configurationPath)
	if err != nil {
		panic(err)
	}

	config := ConfigurationModels.Configuration{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}

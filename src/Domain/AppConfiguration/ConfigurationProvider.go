package AppConfiguration

import (
	"encoding/json"
	"hlsapi/src/Domain/AppConfiguration/Models"
	"os"
)

var configurationPath = ""

type JsonConfigurationProvider struct {
}

func (JsonConfigurationProvider) ReadRoot() ConfigurationModels.ConfigurationRoot {
	if len(configurationPath) == 0 {
		panic("configuration manager was not initialized.")
	}
	configBytes, err := os.ReadFile(configurationPath)
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

func (JsonConfigurationProvider) Initialize(configJsonPath string) {
	configurationPath = configJsonPath
}

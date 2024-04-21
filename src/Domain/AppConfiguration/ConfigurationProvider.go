package AppConfiguration

import (
	"encoding/json"
	"hlsapi/src/Application/Entities"
	"os"
)

var configurationPath = ""

type JsonConfigurationProvider struct {
}

func (JsonConfigurationProvider) GetConfiguration() Entities.AppConfiguration {
	if len(configurationPath) == 0 {
		panic("configuration manager was not initialized.")
	}
	configBytes, err := os.ReadFile(configurationPath)
	if err != nil {
		panic(err)
	}

	config := Entities.AppConfiguration{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func (JsonConfigurationProvider) Initialize(configJsonPath string) {
	configurationPath = configJsonPath
}

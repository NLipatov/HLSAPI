package AppConfiguration

import (
	"encoding/json"
	"hlsapi/src/Application/Entities"
	AppConfigurationErrors "hlsapi/src/Domain/AppConfiguration/Errors"
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

	EnvAppAddress := os.Getenv("APP_ADDRESS")
	if len(EnvAppAddress) > 0 {
		config.Server.AppAddress = EnvAppAddress
	}

	return config
}

func (JsonConfigurationProvider) Initialize(configJsonPath string) error {
	configurationPath = configJsonPath

	//Will update app configuration file with values passed in as env variables
	err := setupEnvValues()
	if err != nil {
		return err
	}
	return nil
}

func setupEnvValues() error {
	if len(configurationPath) == 0 {
		panic("configuration manager was not initialized.")
	}
	configBytes, err := os.ReadFile(configurationPath)
	if err != nil {
		return AppConfigurationErrors.EnvConfigurationUpdateError{
			InnerError: err,
		}
	}

	config := Entities.AppConfiguration{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return AppConfigurationErrors.EnvConfigurationUpdateError{
			InnerError: err,
		}
	}

	EnvAppAddress := os.Getenv("APP_ADDRESS")
	if len(EnvAppAddress) > 0 {
		config.Server.AppAddress = EnvAppAddress
	}

	modifiedConfigBytes, err := json.Marshal(config)
	if err != nil {
		return AppConfigurationErrors.EnvConfigurationUpdateError{
			InnerError: err,
		}
	}

	err = os.WriteFile(configurationPath, modifiedConfigBytes, 0644)
	if err != nil {
		return AppConfigurationErrors.EnvConfigurationUpdateError{
			InnerError: err,
		}
	}

	return nil
}

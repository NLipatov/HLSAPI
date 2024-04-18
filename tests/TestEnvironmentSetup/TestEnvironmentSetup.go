package TestEnvironmentSetup

import (
	"encoding/json"
	"hlsapi/src/Domain/AppConfiguration"
	ConfigurationModels "hlsapi/src/Domain/AppConfiguration/Models"
	"os"
	"path"
)

func SetupTestConfiguration(testTemporaryDirectory string) {
	configurationPath := createConfigurationInTestFolder(testTemporaryDirectory, "testSettings.json")
	AppConfiguration.JsonConfigurationProvider{}.Initialize(configurationPath)
}

func createConfigurationInTestFolder(tempFolderPath string, configurationFilename string) string {
	configurationRoot := ConfigurationModels.ConfigurationRoot{
		Server: ConfigurationModels.ServerConfiguration{
			Port:                   9001,
			GetFileEndpointPostfix: "get?filename=",
		},
		Storage: ConfigurationModels.StorageConfiguration{
			MaxFileSizeMb:     100,
			StorageFolderPath: tempFolderPath,
		},
		StorageDaemon: ConfigurationModels.StorageDaemonConfiguration{},
	}

	jsonBytes, err := json.MarshalIndent(configurationRoot, "", "\t")
	if err != nil {
		panic(err)
	}

	testConfigPath := path.Join(tempFolderPath, configurationFilename)
	f, err := os.OpenFile(testConfigPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.Write(jsonBytes)

	return testConfigPath
}

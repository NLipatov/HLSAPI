package AppConfiguration

import (
	"os"
	"path"
	"testing"
)

func TestSetupEnvValues(t *testing.T) {
	tempDirPath := t.TempDir()
	testConfigurationPath := path.Join(tempDirPath, "appSettings.json")
	configContentString := `{
  "Server": {
    "Port": 9001,
    "AppAddress": "http://localhost:9001",
    "GetFileEndpointPostfix": "get?filename="
  },
  "Storage": {
    "MaxFileSizeMb": 300,
    "StorageFolderPath": "./storage"
  },
  "StorageDaemon": {
    "ShouldRun": true,
    "StorageLimitMinutes": 15,
    "StorageChecksIntervalMinutes": 5,
    "EnableLogging": true
  },
  "InfrastructureLayerConfiguration": {
    "FFMPEGConverter": {
      "UseLogging": true
    }
  }
}`
	err := os.WriteFile(testConfigurationPath, []byte(configContentString), 0777)
	if err != nil {
		panic(err)
	}

	expectedAppAddress := "https://hlsapi.com"
	err = os.Setenv("APP_ADDRESS", expectedAppAddress)
	if err != nil {
		panic(err)
	}
	provider := JsonConfigurationProvider{}
	err = provider.Initialize(testConfigurationPath)
	if err != nil {
		panic(err)
	}

	actualAppAddress := provider.GetConfiguration().Server.AppAddress
	if actualAppAddress != expectedAppAddress {
		t.Errorf("Environment variables are not set in the app configuration\n"+
			"Expected AppAddress: %s, but got: %s", expectedAppAddress, provider.GetConfiguration().Server.AppAddress)
	}
}

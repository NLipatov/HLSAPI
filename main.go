package main

import (
	"errors"
	"fmt"
	"hlsapi/src/Configuration"
	ConfigurationModels "hlsapi/src/Configuration/Models"
	"hlsapi/src/FileEndpoints"
	"hlsapi/src/SentinelServiceDaemon"
	"net/http"
	"os"
)

func main() {
	Configuration.Init(fmt.Sprintf("%s", "appSettings.json"))
	configuration := Configuration.ReadConfiguration()

	createStorageFolder(configuration)

	if configuration.SentinelServiceDaemonConfiguration.ShouldRun {
		go func() {
			fmt.Println("Starting sentinel service daemon")
			SentinelServiceDaemon.Start(configuration.StorageFolderPath)
		}()
	}

	PORT := fmt.Sprintf(":%d", configuration.Port)

	http.HandleFunc("/store", FileEndpoints.StoreFileOnDisk)
	http.HandleFunc("/get", FileEndpoints.GetFileFromDisk)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createStorageFolder(configuration ConfigurationModels.Configuration) {
	_, err := os.Stat(configuration.StorageFolderPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dirCreationError := os.Mkdir(configuration.StorageFolderPath, os.FileMode(744))
			if dirCreationError != nil {
				panic(dirCreationError)
			}
			return
		}
		panic(err)
	}
}

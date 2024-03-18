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

	if configuration.Sentinel.ShouldRun {
		go SentinelServiceDaemon.Start()
	}

	PORT := fmt.Sprintf(":%d", configuration.Server.Port)

	http.HandleFunc("/store", FileEndpoints.StoreFileOnDisk)
	http.HandleFunc("/get", FileEndpoints.GetFileFromDisk)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createStorageFolder(configuration ConfigurationModels.ConfigurationRoot) {
	_, err := os.Stat(configuration.Storage.StorageFolderPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			dirCreationError := os.Mkdir(configuration.Storage.StorageFolderPath, os.FileMode(0700))
			if dirCreationError != nil {
				panic(dirCreationError)
			}
			return
		}
		panic(err)
	}
}

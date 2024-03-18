package SentinelServiceDaemon

import (
	"fmt"
	"hlsapi/src/Configuration"
	ConfigurationModels "hlsapi/src/Configuration/Models"
	"os"
	"path/filepath"
	"time"
)

func Start() {
	for getConfiguration().Sentinel.ShouldRun {
		log("Checking storage folder...")
		processDirectory(getConfiguration().Storage.StorageFolderPath)
		log("Sleep")
		interval := getConfiguration().Sentinel.StorageChecksIntervalMinutes
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func processDirectory(path string) {
	log(fmt.Sprintf("checking %s\n", path))
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if len(entries) == 0 {
		err = os.Remove(path)
		if err != nil {
			panic(err)
		}
		return
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		if entry.Type().IsDir() {
			processDirectory(entryPath)
		} else {
			processFile(entryPath)
		}
	}
}

func processFile(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log(fmt.Sprintf("Could not get file information. Filepath: %s", path))
	}

	modtime := fileInfo.ModTime()
	log(fmt.Sprintf("modtime %d\n", modtime))
	storageExpiresAt := modtime.Add(time.Duration(getConfiguration().Sentinel.StorageLimitMinutes) * time.Minute)
	log(fmt.Sprintf("storageExpiresAt %d\n", storageExpiresAt))
	storageExpired := time.Now().After(storageExpiresAt)
	if storageExpired {
		log("Time exceeded: " + path)
	} else {
		log("Time not exceeded: " + path)
	}
	if storageExpired {
		err = os.Remove(path)
		if err != nil {
			log(fmt.Sprintf("Could not delete a file (%s): %s", err.Error(), path))
			return
		}
		log(fmt.Sprintf("Deleted: %s", path))
	}
}

func log(message string) {
	fmt.Println("Sentinel: ", message)
}

func getConfiguration() ConfigurationModels.ConfigurationRoot {
	configuration := Configuration.ReadConfiguration()
	return configuration
}

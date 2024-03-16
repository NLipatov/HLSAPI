package SentinelServiceDaemon

import (
	"fmt"
	"hlsapi/src/Configuration"
	ConfigurationModels "hlsapi/src/Configuration/Models"
	"os"
	"path/filepath"
	"time"
)

func Start(path string) {
	configuration := getConfiguration()
	if !configuration.SentinelServiceDaemonConfiguration.ShouldRun {
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		log(fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	for _, entry := range entries {
		if entry.Type().IsDir() {
			Start(filepath.Join(path, entry.Name()))
		} else {
			entryPath := filepath.Join(path, entry.Name())
			fileInfo, err := os.Stat(entryPath)
			if err != nil {
				log(fmt.Sprintf("Could not get file information. Filepath: %s", entryPath))
			}

			if int(time.Since(fileInfo.ModTime()).Minutes()) > configuration.SentinelServiceDaemonConfiguration.StorageLimitMinutes {
				err = os.Remove(entryPath)
				if err != nil {
					log(fmt.Sprintf("Error on delition (%s): %s", err.Error(), entryPath))
					continue
				}
				log(fmt.Sprintf("Deleted: %s", entryPath))
			}
		}
	}

	interval := getConfiguration().SentinelServiceDaemonConfiguration.StorageChecksIntervalMinutes
	time.Sleep(time.Duration(interval) * time.Minute)
	Start(configuration.StorageFolderPath)
}

func getConfiguration() ConfigurationModels.Configuration {
	configuration := Configuration.ReadConfiguration()
	return configuration
}

func log(message string) {
	fmt.Println("Sentinel: ", message)
}

package Subdomain

import (
	"fmt"
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"path/filepath"
	"time"
)

func Start() {
	for {
		interval := AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().StorageDaemon.StorageChecksIntervalMinutes
		shouldRun := AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().StorageDaemon.ShouldRun
		if shouldRun {
			time.Sleep(time.Duration(interval) * time.Minute)
		}

		log("Checking storage folder...")
		processDirectory(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath)
		log("Sleep")
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func processDirectory(path string) {
	log(fmt.Sprintf("checking %s\n", path))
	entries, err := os.ReadDir(path)
	if err != nil {
		log(fmt.Sprintf("Could not read path: %s\n", path))
		return
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
	storageExpiresAt := modtime.Add(time.Duration(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().StorageDaemon.StorageLimitMinutes) * time.Minute)
	storageExpired := time.Now().After(storageExpiresAt)

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
	shouldLog := AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().StorageDaemon.EnableLogging
	if shouldLog {
		fmt.Println("Storage Daemon: ", message)
	}
}

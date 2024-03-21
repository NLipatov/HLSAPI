package Domain

import (
	"errors"
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]bool{
	".m3u8": true,
	".ts":   true,
	".m4a":  true,
}

func CanFileBeStored(filename string) bool {
	ext := filepath.Ext(filename)
	return allowedExtensions[ext]
}

func SplitIntoFolderAndFilename(originalFilename string) (string, string) {
	pathSequence := strings.Split(originalFilename, "_")
	folder := pathSequence[0]
	filename := pathSequence[1]

	createFolderIfNotExists(AppConfiguration.ReadRoot().Storage.StorageFolderPath)
	createFolderIfNotExists(filepath.Join(AppConfiguration.ReadRoot().Storage.StorageFolderPath, folder))

	return folder, filename
}

func createFolderIfNotExists(path string) {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			folderCreationError := os.Mkdir(path, 0700)
			if folderCreationError != nil {
				panic(folderCreationError)
			}
		} else {
			panic(err)
		}
	}
}

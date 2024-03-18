package FileEndpoints

import (
	"errors"
	"hlsapi/src/Configuration"
	"os"
	"path/filepath"
	"strings"
)

func StorageFolderAndFilename(originalFilename string) (string, string) {
	pathSequence := strings.Split(originalFilename, "_")
	folder := pathSequence[0]
	filename := pathSequence[1]

	createFolderIfNotExists(Configuration.ReadConfiguration().Storage.StorageFolderPath)
	createFolderIfNotExists(filepath.Join(Configuration.ReadConfiguration().Storage.StorageFolderPath, folder))

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

package Domain

import (
	"errors"
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"path/filepath"
	"strings"
)

func CanFileBeStored(filename string) bool {
	isM3U8 := filepath.Ext(filename) == ".m3u8"
	isTs := filepath.Ext(filename) == ".ts"
	isM4a := filepath.Ext(filename) == ".m4a"

	return isM3U8 || isTs || isM4a
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

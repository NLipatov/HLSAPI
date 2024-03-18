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

	_, err := os.Stat(filepath.Join(Configuration.ReadConfiguration().StorageFolderPath, folder))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			folderCreationError := os.Mkdir(filepath.Join(Configuration.ReadConfiguration().StorageFolderPath, folder), 600)
			if folderCreationError != nil {
				panic(folderCreationError)
			}
		} else {
			panic(err)
		}
	}

	return folder, filename
}

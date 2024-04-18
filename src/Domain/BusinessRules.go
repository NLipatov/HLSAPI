package Domain

import (
	"errors"
	"fmt"
	"github.com/u2takey/go-utils/uuid"
	"hlsapi/src/Domain/AppConfiguration"
	ConfigurationModels "hlsapi/src/Domain/AppConfiguration/Models"
	"os"
	"path/filepath"
	"strings"
)

type ConfigurationProvider interface {
	ReadRoot() ConfigurationModels.ConfigurationRoot
}

var allowedExtensions = map[string]bool{
	".m3u8": true,
	".ts":   true,
	".m4a":  true,
}

var allowedExtensionsForM3U8Conversion = map[string]bool{
	".mp4":   true,
	".mts":   true,
	".avchd": true,
	".3GP":   true,
	".mpg":   true,
	".flv":   true,
	".mkv":   true,
	".wmv":   true,
	".mov":   true,
	".avi":   true,
	".webm":  true,
	".h264":  true,
	".hevc":  true,
}

func CanFileBeConvertedToM3U8(filename string) bool {
	ext := filepath.Ext(filename)
	return allowedExtensionsForM3U8Conversion[ext]
}

func CanFileBeStored(filename string) bool {
	ext := filepath.Ext(filename)
	return allowedExtensions[ext]
}

func GetSequenceStorageFolderAndFilename(originalFilename string) (string, string) {
	pathSequence := strings.Split(originalFilename, "_")
	folder := pathSequence[0]
	filename := pathSequence[1]

	CreateFolder(filepath.Join(AppConfiguration.JsonConfigurationProvider{}.ReadRoot().Storage.StorageFolderPath, folder))

	return folder, filename
}

func CreateWorkdir() string {
	folder := uuid.NewUUID()

	CreateFolder(filepath.Join(AppConfiguration.JsonConfigurationProvider{}.ReadRoot().Storage.StorageFolderPath, folder))

	return folder
}

func CreateFolder(path string) {
	aggregatedSegmentPath := ""
	for _, segment := range strings.Split(path, string(os.PathSeparator)) {
		if len(segment) == 0 {
			panic(fmt.Sprintf("Invalid path: %s\n", path))
		}

		if len(aggregatedSegmentPath) != 0 {
			aggregatedSegmentPath += string(os.PathSeparator)
		}
		aggregatedSegmentPath += segment

		_, err := os.Stat(aggregatedSegmentPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				folderCreationError := os.Mkdir(aggregatedSegmentPath, 0700)
				if folderCreationError != nil {
					panic(folderCreationError)
				}
			} else {
				panic(err)
			}
		}
	}
}

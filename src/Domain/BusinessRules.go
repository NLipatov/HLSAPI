package Domain

import (
	"errors"
	"fmt"
	"github.com/u2takey/go-utils/uuid"
	"hlsapi/src/Application/Boundaries"
	ConfigurationModels "hlsapi/src/Application/Entities"
	"hlsapi/src/Domain/AppConfiguration"
	"hlsapi/src/Domain/CleanupType"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ConfigurationProvider interface {
	GetConfiguration() ConfigurationModels.AppConfiguration
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

	CreateFolder(filepath.Join(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, folder))

	return folder, filename
}

func CreateWorkdir() string {
	folder := uuid.NewUUID()

	CreateFolder(filepath.Join(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, folder))

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

func ShouldFileBeCleanedUp(filepath string, mode CleanupType.CleanupMode, environmentManager Boundaries.EnvironmentBoundary, configurationManager Boundaries.ConfigurationBoundary) bool {
	baseWorkdir := path.Join(environmentManager.GetAppRootPath(), configurationManager.GetConfiguration().Storage.StorageFolderPath)

	//File is not in storage folder
	if !strings.HasPrefix(filepath, baseWorkdir) {
		return false
	}

	//File does not exist
	_, err := os.Stat(filepath)
	if err != nil {
		return false
	}

	switch mode {
	case CleanupType.UNSET:
		return false
	case CleanupType.REMOVE_ALL_FILES:
		return true
	case CleanupType.REMOVE_KEY_FILES:
		if strings.HasSuffix(filepath, ".key") || strings.HasSuffix(filepath, ".keyinfo") {
			return true
		}
	default:
		panic("Invalid cleanup mode")
	}

	return false
}

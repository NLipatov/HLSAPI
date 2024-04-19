package Application

import (
	"errors"
	"fmt"
	"hlsapi/src/Application/Boundaries"
	"hlsapi/src/Application/Errors"
	"hlsapi/src/Domain"
	"hlsapi/src/Domain/AppConfiguration"
	"hlsapi/src/Domain/CleanupType"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func StoreFileOnDisk(filename string, readFrom io.Reader, boundary Boundaries.StoreBoundary) error {
	if !Domain.CanFileBeConvertedToM3U8(filename) {
		return ApplicationLayerErrors.FileCantBeStored{}
	}

	folder, filename := Domain.GetSequenceStorageFolderAndFilename(filename)
	path := filepath.Join(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, folder, filename)
	writeTo, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer writeTo.Close()

	err = boundary.Store(writeTo, readFrom)
	if err != nil {
		return err
	}

	return nil
}

func GetFileFromDisk(writeTo io.Writer, requestedFileCode string, boundary Boundaries.GetBoundary) error {
	folder, filename := Domain.GetSequenceStorageFolderAndFilename(requestedFileCode)
	path := filepath.Join(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, folder, filename)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return err
	}

	err := boundary.Get(writeTo, path)
	if err != nil {
		return err
	}

	return nil
}

func ConvertVideoToM3U8Playlist(filename string, readFrom io.Reader, mediaConverterBoundary Boundaries.MediaConverterBoundary, envBoundary Boundaries.EnvironmentBoundary) (playlistContent string, err error) {
	if !Domain.CanFileBeConvertedToM3U8(filename) {
		return "", ApplicationLayerErrors.FileCantBeConvertedToM3U8{}
	}

	workdir := Domain.CreateWorkdir()
	outputFolderPath := filepath.Join(AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, workdir)

	inputFilePath := filepath.Join(outputFolderPath, filename)
	newFileName := "in" + path.Ext(inputFilePath)

	writeTo, err := os.Create(filepath.Join(outputFolderPath, newFileName))
	if err != nil {
		panic(err)
	}
	defer writeTo.Close()

	_, err = io.Copy(writeTo, readFrom)
	if err != nil {
		panic(err)
	}

	appRoot := envBoundary.GetAppRootPath()
	playlistPath, err := mediaConverterBoundary.ConvertToM3U8(path.Join(appRoot, "storage", workdir), "in"+path.Ext(inputFilePath))
	if err != nil {
		panic(err)
	}

	m3u8Data, err := os.ReadFile(playlistPath)
	if err != nil {
		return "", err
	}

	formattedPlaylist, err := formatPlaylist(string(m3u8Data), workdir)
	if err != nil {
		panic(err)
	}

	return formattedPlaylist, nil

}

func formatPlaylist(playlist string, folderId string) (string, error) {
	port := AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Server.Port
	endpointPostfix := AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Server.GetFileEndpointPostfix
	url := fmt.Sprintf("http://localhost:%d/%s%s_", port, endpointPostfix, folderId)

	sb := strings.Builder{}

	lines := strings.Split(playlist, "\n")
	for _, line := range lines {
		formattedLine := strings.TrimSpace(line)
		formattedLine = strings.ReplaceAll(formattedLine, "\n", "")
		formattedLine = strings.ReplaceAll(formattedLine, "\t", "")

		if strings.HasPrefix(formattedLine, "#EXT-X-KEY") && strings.Contains(formattedLine, "URI=\"./file.key\"") {
			formattedLine = strings.Replace(formattedLine, "./file.key", fmt.Sprintf("%sfile.key", url), 1)
		} else if strings.HasPrefix(formattedLine, "out") {
			formattedLine = fmt.Sprintf("%s%s", url, formattedLine)
		} else {
			formattedLine = fmt.Sprintf("%s", formattedLine)
		}

		sb.WriteString(fmt.Sprintf("%s\n", formattedLine))
	}

	return sb.String(), nil
}

func CleanUpStorageFolder(folderId string, mode CleanupType.CleanupMode, environmentManager Boundaries.EnvironmentBoundary) error {
	folder := filepath.Join(environmentManager.GetAppRootPath(), AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, folderId)

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := file.Name()
		if Domain.ShouldFileBeCleanedUp(path.Join(folder, filename), mode, environmentManager, AppConfiguration.JsonConfigurationProvider{}) {
			filePath := path.Join(folder, filename)
			err = os.Remove(filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

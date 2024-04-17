package Application

import (
	"errors"
	"hlsapi/src/Application/Boundaries"
	"hlsapi/src/Application/Errors"
	"hlsapi/src/Domain"
	"hlsapi/src/Domain/AppConfiguration"
	"io"
	"os"
	"path/filepath"
)

func StoreFileOnDisk(filename string, readFrom io.Reader, boundary Boundaries.StoreBoundary) error {
	if !Domain.CanFileBeConvertedToM3U8(filename) {
		return ApplicationLayerErrors.FileCantBeStored{}
	}

	folder, filename := Domain.GetSequenceStorageFolderAndFilename(filename)
	path := filepath.Join(AppConfiguration.JsonConfigurationProvider{}.ReadRoot().Storage.StorageFolderPath, folder, filename)
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
	path := filepath.Join(AppConfiguration.JsonConfigurationProvider{}.ReadRoot().Storage.StorageFolderPath, folder, filename)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return err
	}

	err := boundary.Get(writeTo, path)
	if err != nil {
		return err
	}

	return nil
}

func ConvertToM3U8(filename string, readFrom io.Reader, boundary Boundaries.MediaConverterBoundary) (string, error) {
	if !Domain.CanFileBeConvertedToM3U8(filename) {
		return "", ApplicationLayerErrors.FileCantBeConvertedToM3U8{}
	}

	folderId := Domain.PrepareFolder()
	outputFolderPath := filepath.Join(AppConfiguration.JsonConfigurationProvider{}.ReadRoot().Storage.StorageFolderPath, folderId)

	inputFilePath := filepath.Join(outputFolderPath, filename)
	outputFilePath := filepath.Join(outputFolderPath, "$.m3u8")

	writeTo, err := os.Create(inputFilePath)
	if err != nil {
		panic(err)
	}
	defer writeTo.Close()

	_, err = io.Copy(writeTo, readFrom)
	if err != nil {
		panic(err)
	}

	err = boundary.ConvertToM3U8(inputFilePath, outputFilePath)
	if err != nil {
		panic(err)
	}

	m3u8Data, err := os.ReadFile(outputFilePath)
	if err != nil {
		return "", err
	}

	return string(m3u8Data), nil
}

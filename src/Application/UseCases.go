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

func StoreFileOnDisk(filename string, src io.Reader, boundary Boundaries.StoreBoundary) error {
	if !Domain.CanFileBeStored(filename) {
		return ApplicationLayerErrors.FileCantBeStored{}
	}

	folder, filename := Domain.SplitIntoFolderAndFilename(filename)
	path := filepath.Join(AppConfiguration.ReadRoot().Storage.StorageFolderPath, folder, filename)
	writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)

	defer writer.Close()
	if err != nil {
		return err
	}

	err = boundary.Store(writer, src)
	if err != nil {
		return err
	}

	return nil
}

func GetFileFromDisk(writeTo io.Writer, requestedFileCode string, boundary Boundaries.GetBoundary) error {
	folder, filename := Domain.SplitIntoFolderAndFilename(requestedFileCode)
	path := filepath.Join(AppConfiguration.ReadRoot().Storage.StorageFolderPath, folder, filename)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return err
	}

	err := boundary.Get(writeTo, path)
	if err != nil {
		return err
	}

	return nil
}

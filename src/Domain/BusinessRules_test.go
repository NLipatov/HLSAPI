package Domain

import (
	ConfigurationModels "hlsapi/src/Domain/AppConfiguration/Models"
	"os"
	"strings"
	"testing"
)

func TestCanFileBeStored_mp3(t *testing.T) {
	filename := "sample.mp3"
	expected := false
	actual := CanFileBeStored(filename)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestCanFileBeStored_ts(t *testing.T) {
	filename := "sample.ts"
	expected := true
	actual := CanFileBeStored(filename)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestCanFileBeStored_m3u8(t *testing.T) {
	filename := "sample.m3u8"
	expected := true
	actual := CanFileBeStored(filename)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

func TestCanFileBeStored_m4a(t *testing.T) {
	filename := "sample.m4a"
	expected := true
	actual := CanFileBeStored(filename)

	if actual != expected {
		t.Errorf("Expected: %v, got: %v", expected, actual)
	}
}

type mockConfigProvider struct {
	MaxFileSizeMb     int
	StorageFolderPath string
}

func (mCP mockConfigProvider) ReadRoot() ConfigurationModels.ConfigurationRoot {
	return ConfigurationModels.ConfigurationRoot{
		Server: ConfigurationModels.ServerConfiguration{},
		Storage: ConfigurationModels.StorageConfiguration{
			MaxFileSizeMb:     mCP.MaxFileSizeMb,
			StorageFolderPath: mCP.StorageFolderPath,
		},
		StorageDaemon: ConfigurationModels.StorageDaemonConfiguration{},
	}
}

func TestGetStorageFolderAndFilename_validFolderAndValidFilename(t *testing.T) {
	mCP := mockConfigProvider{
		MaxFileSizeMb:     0,
		StorageFolderPath: t.TempDir(),
	}
	originalFilename := "folder_filename.ts"
	expectedFolder, expectedFilename := "folder", "filename.ts"
	actualFolder, actualFilename := GetStorageFolderAndFilename(originalFilename, mCP)

	if actualFolder != expectedFolder || actualFilename != expectedFilename {
		t.Errorf("Expected: %v and %v, got: %v and %v", expectedFolder, expectedFilename, actualFolder, actualFilename)
	}
}

func TestGetStorageFolderAndFilename_folderIsActuallyCreated(t *testing.T) {
	mCP := mockConfigProvider{
		MaxFileSizeMb:     0,
		StorageFolderPath: t.TempDir(),
	}
	originalFilename := "folder_filename.ts"
	expectedFolder, expectedFilename := "folder", "filename.ts"
	actualFolder, actualFilename := GetStorageFolderAndFilename(originalFilename, mCP)

	resultingFolderPath := strings.Join([]string{mCP.StorageFolderPath, expectedFolder}, string(os.PathSeparator))
	_, err := os.Stat(resultingFolderPath)

	if err != nil {
		t.Errorf("Expected: folder created: %v, got: folder was not created", resultingFolderPath)
	}

	if actualFolder != expectedFolder || actualFilename != expectedFilename {
		t.Errorf("Expected: %v and %v, got: %v and %v", expectedFolder, expectedFilename, actualFolder, actualFilename)
	}
}

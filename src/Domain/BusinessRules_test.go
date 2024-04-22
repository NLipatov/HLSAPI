package Domain

import (
	"fmt"
	"github.com/google/uuid"
	"hlsapi/src/Domain/AppConfiguration"
	"hlsapi/src/Domain/WipeModes"
	"hlsapi/src/Infrastructure"
	"hlsapi/tests/TestEnvironmentSetup"
	"os"
	"path"
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

func TestGetStorageFolderAndFilename_validFolderAndValidFilename(t *testing.T) {
	TestEnvironmentSetup.SetupTestDirConfiguration(t.TempDir())
	originalFilename := "folder_filename.ts"
	expectedFolder, expectedFilename := "folder", "filename.ts"
	actualFolder, actualFilename := GetSequenceStorageFolderAndFilename(originalFilename)

	if actualFolder != expectedFolder || actualFilename != expectedFilename {
		t.Errorf("Expected: %v and %v, got: %v and %v", expectedFolder, expectedFilename, actualFolder, actualFilename)
	}
}

func TestGetStorageFolderAndFilename_folderIsActuallyCreated(t *testing.T) {
	TestEnvironmentSetup.SetupTestDirConfiguration(t.TempDir())
	originalFilename := "folder_filename.ts"
	expectedFolder, expectedFilename := "folder", "filename.ts"
	actualFolder, actualFilename := GetSequenceStorageFolderAndFilename(originalFilename)

	resultingFolderPath := strings.Join([]string{AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Storage.StorageFolderPath, expectedFolder}, string(os.PathSeparator))
	_, err := os.Stat(resultingFolderPath)

	if err != nil {
		t.Errorf("Expected: folder created: %v, got: folder was not created", resultingFolderPath)
	}

	if actualFolder != expectedFolder || actualFilename != expectedFilename {
		t.Errorf("Expected: %v and %v, got: %v and %v", expectedFolder, expectedFilename, actualFolder, actualFilename)
	}
}

func TestShouldFileBeCleanedUp(t *testing.T) {
	tempDir := t.TempDir()
	TestEnvironmentSetup.SetupTestDirConfiguration(tempDir)

	keyPath := path.Join(tempDir, "key.key")
	keyinfoPath := path.Join(tempDir, "key.keyinfo")
	m3u8Path := path.Join(tempDir, "sample.m3u8")
	tsPath := path.Join(tempDir, "sample.ts")
	notExistingFilePath := path.Join(tempDir, fmt.Sprintf("%s.%s", uuid.New().String(), uuid.New().String()))
	fileFromOtherDir := path.Join("other", fmt.Sprintf("%s.%s", uuid.New().String(), uuid.New().String()))

	_ = os.WriteFile(keyPath, []byte("some content"), 0777)
	_ = os.WriteFile(keyinfoPath, []byte("some content"), 0777)
	_ = os.WriteFile(m3u8Path, []byte("some content"), 0777)
	_ = os.WriteFile(tsPath, []byte("some content"), 0777)

	keyFileResult := ShouldFileBeCleanedUp(keyPath, WipeModes.REMOVE_KEY_FILES, Infrastructure.EnvironmentManager{}, AppConfiguration.JsonConfigurationProvider{})
	keyinfoFileResult := ShouldFileBeCleanedUp(keyinfoPath, WipeModes.REMOVE_KEY_FILES, Infrastructure.EnvironmentManager{}, AppConfiguration.JsonConfigurationProvider{})
	m3u8FileResult := ShouldFileBeCleanedUp(m3u8Path, WipeModes.REMOVE_KEY_FILES, Infrastructure.EnvironmentManager{}, AppConfiguration.JsonConfigurationProvider{})
	tsFileResult := ShouldFileBeCleanedUp(tsPath, WipeModes.REMOVE_KEY_FILES, Infrastructure.EnvironmentManager{}, AppConfiguration.JsonConfigurationProvider{})
	notExistingFilePathResult := ShouldFileBeCleanedUp(notExistingFilePath, WipeModes.REMOVE_KEY_FILES, Infrastructure.EnvironmentManager{}, AppConfiguration.JsonConfigurationProvider{})
	fileFromOtherDirResult := ShouldFileBeCleanedUp(fileFromOtherDir, WipeModes.REMOVE_ALL_FILES, Infrastructure.EnvironmentManager{}, AppConfiguration.JsonConfigurationProvider{})

	if keyFileResult != true {
		t.Error(".key file should be cleaned up")
	}

	if keyinfoFileResult != true {
		t.Error(".keyinfo file should be cleaned up")
	}

	if m3u8FileResult != false {
		t.Error(".m3u8 file should not be cleaned up")
	}

	if tsFileResult != false {
		t.Error(".ts file should not be cleaned up")
	}

	if notExistingFilePathResult != false {
		t.Error("Not existing file should not be cleaned up")
	}

	if fileFromOtherDirResult == true {
		t.Errorf("File from other dir should not be cleaned up")
	}
}

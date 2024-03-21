package tests

import (
	"bufio"
	"bytes"
	"hlsapi/src/Application"
	"hlsapi/src/Domain/AppConfiguration"
	"hlsapi/src/Infrastructure"
	"hlsapi/tests/TestEnvironmentSetup"
	"os"
	"strings"
	"testing"
)

func TestGetFileFromDisk(t *testing.T) {
	testDir := t.TempDir()
	testFileContent := []byte{132, 243, 0, 73}
	testFilePath := writeTestFile(testDir, testFileContent)
	pathArray := strings.Split(testFilePath, string(os.PathSeparator))
	folder, filename := pathArray[len(pathArray)-2], pathArray[len(pathArray)-1]

	writeTo := bytes.Buffer{}
	err := Application.GetFileFromDisk(&writeTo, strings.Join([]string{folder, filename}, "_"), Infrastructure.DiskInteractor{})
	if err != nil {
		t.Error(err)
	}

	for i, v := range writeTo.Bytes() {
		if testFileContent[i] != v {
			t.Errorf("File content is invalid")
		}
	}
}

func TestStoreFileOnDisk(t *testing.T) {
	testDir := t.TempDir()

	testFileContent := []byte{132, 243, 0, 73}
	storedTestFilePath := writeTestFile(testDir, testFileContent)
	f, err := os.Open(storedTestFilePath)
	if err != nil {
		t.Error(err)
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	buffer := make([]byte, len(testFileContent))
	_, err = reader.Read(buffer)

	if err != nil {
		t.Error(err)
	}

	if len(buffer) == 0 {
		t.Error("File was empty")
	}

	for i, v := range buffer {
		if testFileContent[i] != v {
			t.Error("File was written, but content is invalid")
		}
	}
}

func writeTestFile(tempDirPath string, content []byte) string {
	setupTestConfiguration(tempDirPath)
	testFolder := "test"
	testFilename := "sample.m3u8"
	testReader := bytes.NewReader(content)
	err := Application.StoreFileOnDisk(strings.Join([]string{testFolder, testFilename}, "_"), testReader, Infrastructure.DiskInteractor{})

	if err != nil {
		panic(err)
	}

	storedTestFilePath := strings.Join([]string{tempDirPath, testFolder, testFilename}, string(os.PathSeparator))
	return storedTestFilePath
}

func setupTestConfiguration(testTemporaryDirectory string) {
	configurationPath := TestEnvironmentSetup.CreateConfigurationInTestFolder(testTemporaryDirectory, "testSettings.json")
	AppConfiguration.JsonConfigurationProvider{}.Initialize(configurationPath)
}

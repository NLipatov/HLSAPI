package Application

import (
	"bufio"
	"bytes"
	"encoding/json"
	"hlsapi/src/Domain/AppConfiguration"
	ConfigurationModels "hlsapi/src/Domain/AppConfiguration/Models"
	"hlsapi/src/Infrastructure"
	"os"
	"path"
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
	err := GetFileFromDisk(&writeTo, strings.Join([]string{folder, filename}, "_"), Infrastructure.DiskInteractor{})
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
	err := StoreFileOnDisk(strings.Join([]string{testFolder, testFilename}, "_"), testReader, Infrastructure.DiskInteractor{})

	if err != nil {
		panic(err)
	}

	storedTestFilePath := strings.Join([]string{tempDirPath, testFolder, testFilename}, string(os.PathSeparator))
	return storedTestFilePath
}

func setupTestConfiguration(testTemporaryDirectory string) {
	configurationPath := createConfigurationInTestFolder(testTemporaryDirectory, "testSettings.json")
	AppConfiguration.JsonConfigurationProvider{}.Initialize(configurationPath)
}

func createConfigurationInTestFolder(tempFolderPath string, configurationFilename string) string {
	configurationRoot := ConfigurationModels.ConfigurationRoot{
		Server: ConfigurationModels.ServerConfiguration{},
		Storage: ConfigurationModels.StorageConfiguration{
			MaxFileSizeMb:     100,
			StorageFolderPath: tempFolderPath,
		},
		StorageDaemon: ConfigurationModels.StorageDaemonConfiguration{},
	}

	jsonBytes, err := json.MarshalIndent(configurationRoot, "", "\t")
	if err != nil {
		panic(err)
	}

	testConfigPath := path.Join(tempFolderPath, configurationFilename)
	f, err := os.OpenFile(testConfigPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	f.Write(jsonBytes)

	return testConfigPath
}

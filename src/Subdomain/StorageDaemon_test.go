package Subdomain

import (
	"fmt"
	"github.com/google/uuid"
	"hlsapi/tests/TestEnvironmentSetup"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessDirectoryWithOutdatedFiles(t *testing.T) {
	tempDir := os.TempDir()
	TestEnvironmentSetup.SetupTestDirConfiguration(tempDir)

	testDirPath := filepath.Join(tempDir, uuid.New().String())
	err := os.Mkdir(testDirPath, 0700)
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(testDirPath)

	//in testDirPath creates 3 folders (0, 1, 2), each of them contains 3 files(0.tmp, 1.tmp, 2.tmp)
	for i := 0; i < 3; i++ {
		folderPath := filepath.Join(testDirPath, fmt.Sprintf("%d", i))
		err = os.Mkdir(folderPath, 0700)
		if err != nil {
			t.Error(err)
		}

		for j := 0; j < 3; j++ {
			filePath := filepath.Join(folderPath, fmt.Sprintf("%d.tmp", j))
			file, err := os.Create(filePath)
			if err != nil {
				t.Error(err)
			}

			file.Close()

			err = os.Chtimes(filePath, time.Now(), time.Now().Add(time.Minute*-10))
			if err != nil {
				t.Error(err)
			}
		}
	}

	processDirectory(testDirPath)

	_, err = os.Stat(testDirPath)
	if os.IsNotExist(err) {
		//Ok, folder with outdated files deleted
		return
	}

	//Storage daemon failed to clean up folder with outdated files
	t.Error(fmt.Sprintf("Directory with outdated files exists after processing: %s", testDirPath))
}

func TestProcessDirectoryWithNotOutdatedFiles(t *testing.T) {
	tempDir := os.TempDir()
	TestEnvironmentSetup.SetupTestDirConfiguration(tempDir)

	testDirPath := filepath.Join(tempDir, uuid.New().String())
	err := os.Mkdir(testDirPath, 0700)
	if err != nil {
		t.Error(err)
	}

	defer os.RemoveAll(testDirPath)

	//in testDirPath creates 3 folders (0, 1, 2), each of them contains 3 files(0.tmp, 1.tmp, 2.tmp)
	for i := 0; i < 3; i++ {
		folderPath := filepath.Join(testDirPath, fmt.Sprintf("%d", i))
		err = os.Mkdir(folderPath, 0700)
		if err != nil {
			t.Error(err)
		}

		for j := 0; j < 3; j++ {
			filePath := filepath.Join(folderPath, fmt.Sprintf("%d.tmp", j))
			file, err := os.Create(filePath)
			if err != nil {
				t.Error(err)
			}

			file.Close()
		}
	}

	processDirectory(testDirPath)

	_, err = os.Stat(testDirPath)
	if os.IsNotExist(err) {
		t.Error(fmt.Sprintf("Directory should not be removed after processing: %s", testDirPath))
	}

	//Checking that all the files are still there
	for i := 0; i < 3; i++ {
		folderPath := filepath.Join(testDirPath, fmt.Sprintf("%d", i))

		for j := 0; j < 3; j++ {
			filePath := filepath.Join(folderPath, fmt.Sprintf("%d.tmp", j))
			_, err := os.Stat(filePath)
			if err != nil {
				t.Error(err)
			}
			if os.IsNotExist(err) {
				t.Error(fmt.Sprintf("File should not be removed after processing: %s", filePath))
			}
		}
	}
}

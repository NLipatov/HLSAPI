package Application

import (
	"bufio"
	"bytes"
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

func TestFormatPlaylist(t *testing.T) {
	tempDir := t.TempDir()
	TestEnvironmentSetup.SetupTestDirConfiguration(tempDir)

	playlist :=
		`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:14
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-PLAYLIST-TYPE:VOD
#EXT-X-KEY:METHOD=AES-128,URI="./file.key",IV=0xfd0c8974ce8f67645add2064c3fd3104
#EXTINF:10.463000,
out0.ts
#EXTINF:12.625000,
out1.ts
#EXTINF:8.958000,
out2.ts
#EXTINF:10.417000,
out3.ts
#EXTINF:13.666000,
out4.ts
#EXTINF:3.854000,
out5.ts
#EXT-X-ENDLIST`

	formattedPlaylist, err := formatPlaylist(playlist, "4c1cb6c6-b00c-4c07-8954-3332cc668c83")
	if err != nil {
		t.Error(err)
	}

	expectedPlaylist :=
		`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:14
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-PLAYLIST-TYPE:VOD
#EXT-X-KEY:METHOD=AES-128,URI="https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_file.key",IV=0xfd0c8974ce8f67645add2064c3fd3104
#EXTINF:10.463000,
https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_out0.ts
#EXTINF:12.625000,
https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_out1.ts
#EXTINF:8.958000,
https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_out2.ts
#EXTINF:10.417000,
https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_out3.ts
#EXTINF:13.666000,
https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_out4.ts
#EXTINF:3.854000,
https://example.com/get?filename=4c1cb6c6-b00c-4c07-8954-3332cc668c83_out5.ts
#EXT-X-ENDLIST`

	if strings.TrimSpace(expectedPlaylist) != strings.TrimSpace(formattedPlaylist) {
		t.Error("Playlist is invalid")
	}
}

func writeTestFile(tempDirPath string, content []byte) string {
	TestEnvironmentSetup.SetupTestDirConfiguration(tempDirPath)
	testFolder := "test"
	testFilename := "sample.mp4"
	testReader := bytes.NewReader(content)
	err := StoreFileOnDisk(strings.Join([]string{testFolder, testFilename}, "_"), testReader, Infrastructure.DiskInteractor{})

	if err != nil {
		panic(err)
	}

	storedTestFilePath := strings.Join([]string{tempDirPath, testFolder, testFilename}, string(os.PathSeparator))
	return storedTestFilePath
}

package Infrastructure

import (
	"fmt"
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"os/exec"
	"path"
)

type FfmpegConverter struct{}

func (f FfmpegConverter) ConvertToM3U8(inputFilename string, outputFilename string) error {
	dirPrefix := "/app/" + "storage/" + path.Base(inputFilename) + "/"
	fmt.Println("dirPrefix: %s", dirPrefix)
	script := "#!/bin/bash\nBASE_URL=${1:-'.'}\nopenssl rand 16 > " + dirPrefix + "file.key" + "\necho $BASE_URL/file.key > " + dirPrefix + "file.keyinfo" + "\necho " + dirPrefix + "file.key" + " >> " + dirPrefix + "file.keyinfo" + "\necho $(openssl rand -hex 16) >> " + dirPrefix + "file.keyinfo" + "\nffmpeg -i " + dirPrefix + "frag_bunny.avi" + " -c:v copy -c:a copy -hls_time 10 -hls_playlist_type vod -hls_key_info_file " + dirPrefix + "file.keyinfo " + dirPrefix + "out.m3u8"
	scriptPath := dirPrefix + "hls.sh"

	// Записываем скрипт в файл
	err := os.WriteFile(scriptPath, []byte(script), 0777)
	if err != nil {
		return err
	}

	// Делаем скрипт исполняемым
	cmd := exec.Command("chmod", "+x", scriptPath)
	if isLoggingEnabled() {
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Запускаем скрипт
	cmd = exec.Command(scriptPath)
	if isLoggingEnabled() {
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func isLoggingEnabled() bool {
	return AppConfiguration.JsonConfigurationProvider{}.ReadRoot().InfrastructureLayerConfiguration.FFMPEGConverter.UseLogging
}

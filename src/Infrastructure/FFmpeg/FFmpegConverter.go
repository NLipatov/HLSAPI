package FFmpeg

import (
	"fmt"
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"os/exec"
	"path"
)

type Converter struct{}

func (f Converter) ConvertToM3U8(workdirAbsPath string, inputVideoFilename string) (string, error) {
	scriptPath, err := generateSh(workdirAbsPath, inputVideoFilename)
	if err != nil {
		fmt.Println("Couldn't generate FFmpeg sh script")
		panic(err)
	}

	err = makeShExecutable(scriptPath)
	if err != nil {
		fmt.Println("Couldn't make FFmpeg sh script executable")
		panic(err)
	}

	err = executeSh(scriptPath)
	if err != nil {
		fmt.Println("Couldn't execute FFmpeg sh script")
		panic(err)
	}

	tryRemovePostProcessFiles(path.Join(workdirAbsPath, inputVideoFilename), path.Join(workdirAbsPath, path.Base(scriptPath)))

	playlistPath := path.Join(workdirAbsPath, "out.m3u8")

	return playlistPath, nil
}

func isLoggingEnabled() bool {
	return AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().InfrastructureLayerConfiguration.FFMPEGConverter.UseLogging
}

func tryRemovePostProcessFiles(paths ...string) {
	for _, p := range paths {
		_ = os.Remove(p)
	}
}

func generateSh(AbsoluteFolderPath string, inVideoFilename string) (string, error) {
	script := "#!/bin/bash\nBASE_URL=${1:-'.'}\nopenssl rand 16 > " + path.Join(AbsoluteFolderPath, "file.key") + "\necho $BASE_URL/file.key > " + path.Join(AbsoluteFolderPath, "file.keyinfo") + "\necho " + path.Join(AbsoluteFolderPath, "file.key") + " >> " + path.Join(AbsoluteFolderPath, "file.keyinfo") + "\necho $(openssl rand -hex 16) >> " + path.Join(AbsoluteFolderPath, "file.keyinfo") + "\nffmpeg -i " + path.Join(AbsoluteFolderPath, inVideoFilename) + " -c:v copy -c:a copy -hls_time 10 -hls_playlist_type vod -hls_playlist_type vod -hls_key_info_file " + path.Join(AbsoluteFolderPath, "file.keyinfo") + " " + path.Join(AbsoluteFolderPath, "out.m3u8")
	scriptPath := path.Join(AbsoluteFolderPath, "hls.sh")

	err := os.WriteFile(scriptPath, []byte(script), 0777)
	if err != nil {
		return "", err
	}

	return scriptPath, nil
}

func makeShExecutable(scriptPath string) error {
	cmd := exec.Command("chmod", "+x", scriptPath)
	if isLoggingEnabled() {
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func executeSh(scriptPath string) error {
	cmd := exec.Command(scriptPath)
	if isLoggingEnabled() {
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

package FFmpeg

import (
	"fmt"
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Converter struct{}

func (f Converter) ConvertToM3U8(workdirAbsPath string, inputVideoFilename string) (string, error) {
	log("Generating a hls.sh")
	scriptPath, err := generateSh(workdirAbsPath, inputVideoFilename)
	if err != nil {
		fmt.Println("Couldn't generate hls.sh script")
		panic(err)
	}

	log("Making hls.sh an executable")
	err = makeShExecutable(scriptPath)
	if err != nil {
		fmt.Println("Couldn't make hls.sh script executable")
		panic(err)
	}

	log("Executing hls.sh")
	err = executeSh(scriptPath)
	if err != nil {
		fmt.Println("Couldn't execute hls.sh script")
		panic(err)
	}

	tryRemovePostProcessFiles(path.Join(workdirAbsPath, changeExtensionToMp4(inputVideoFilename)), path.Join(workdirAbsPath, inputVideoFilename), path.Join(workdirAbsPath, path.Base(scriptPath)))

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
	script := "#!/bin/bash" +
		"\nBASE_URL=${1:-'.'}" +
		"\nopenssl rand 16 > " + path.Join(AbsoluteFolderPath, "file.key") +
		"\necho $BASE_URL/file.key > " + path.Join(AbsoluteFolderPath, "file.keyinfo") +
		"\necho " + path.Join(AbsoluteFolderPath, "file.key") + " >> " + path.Join(AbsoluteFolderPath, "file.keyinfo") +
		"\necho $(openssl rand -hex 16) >> " + path.Join(AbsoluteFolderPath, "file.keyinfo")

	if !strings.HasSuffix(strings.ToLower(inVideoFilename), ".mp4") {
		log(fmt.Sprintf("%s will be converted to %s", inVideoFilename, changeExtensionToMp4(inVideoFilename)))
		script += fmt.Sprintf("\nffmpeg -i %s -c:v libx264 -preset ultrafast -c:a aac -b:a 128k %s", path.Join(AbsoluteFolderPath, inVideoFilename), path.Join(AbsoluteFolderPath, changeExtensionToMp4(inVideoFilename)))

		inVideoFilename = changeExtensionToMp4(inVideoFilename)
		log(fmt.Sprintf("New target file is: %s", inVideoFilename))
	}
	script += "\nffmpeg -y -i " + path.Join(AbsoluteFolderPath, inVideoFilename) +
		" -c:v copy" +
		" -c:a copy" +
		" -hls_time 10" +
		" -hls_playlist_type vod" +
		" -threads 0" +
		" -hls_key_info_file " + path.Join(AbsoluteFolderPath, "file.keyinfo") +
		fmt.Sprintf(" %s", path.Join(AbsoluteFolderPath, "out.m3u8"))

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

func log(message string) {
	fmt.Printf("[FFmpegConverter]: %s\n", message)
}

func changeExtensionToMp4(originalPath string) string {
	return strings.TrimSuffix(originalPath, path.Ext(originalPath)) + ".mp4"
}

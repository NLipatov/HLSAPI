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
		log(fmt.Sprintf("%s will be converted to %s before generating a m3u8 playlist", inVideoFilename, changeExtensionToMp4(inVideoFilename)))

		videoCodec, err := getVideoCodec(path.Join(AbsoluteFolderPath, inVideoFilename))
		videoCodecDirectCopyStatus := isDirectVideoCodecCopyPossible(videoCodec)
		log(fmt.Sprintf("%s video codec: %s, can be copied: %t", inVideoFilename, videoCodec, videoCodecDirectCopyStatus))
		VideoCodecArg := "-c:v libx264"
		if err == nil && videoCodecDirectCopyStatus {
			VideoCodecArg = "-c:v copy"
		}

		AudioCodecArg := "-c:a aac"
		audioCodec, err := getAudioCodec(path.Join(AbsoluteFolderPath, inVideoFilename))
		audioCodecDirectCopyStatus := isDirectAudioCodecCopyPossible(audioCodec)
		log(fmt.Sprintf("%s audio codec: %s, can be copied: %t", inVideoFilename, audioCodec, audioCodecDirectCopyStatus))
		if err == nil && audioCodecDirectCopyStatus {
			AudioCodecArg = "-c:a copy"
		}

		convertationToMp4Script := fmt.Sprintf("\nffmpeg -i %s %s -preset ultrafast %s %s", path.Join(AbsoluteFolderPath, inVideoFilename), VideoCodecArg, AudioCodecArg, path.Join(AbsoluteFolderPath, changeExtensionToMp4(inVideoFilename)))
		script += convertationToMp4Script

		log(fmt.Sprintf("Convertation to mp4 script: %s", convertationToMp4Script))

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

func getVideoCodec(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	videoCodecProbeResult := string(out)
	resultParts := strings.Split(strings.TrimSpace(videoCodecProbeResult), "\n")

	if len(resultParts) > 1 {
		reference := resultParts[0]
		for i := 1; i < len(resultParts); i++ {
			if reference != resultParts[i] {
				return videoCodecProbeResult, nil
			}
		}
		return reference, nil
	}

	return strings.TrimSpace(resultParts[0]), nil
}

func getAudioCodec(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", filePath)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func isDirectVideoCodecCopyPossible(codec string) bool {
	allowedCodecs := map[string]bool{
		"h264":       true,
		"hevc":       true,
		"mpeg4":      true,
		"mpeg2video": true,
		"mjpeg":      true,
		"dvvideo":    true,
	}

	return allowedCodecs[codec]
}

func isDirectAudioCodecCopyPossible(codec string) bool {
	allowedCodecs := map[string]bool{
		"aac": true,
		"mp3": true,
	}

	return allowedCodecs[codec]
}

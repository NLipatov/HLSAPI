package Infrastructure

import (
	"hlsapi/src/Domain/AppConfiguration"
	"os"
	"os/exec"
)

type FfmpegConverter struct{}

func (f FfmpegConverter) ConvertToM3U8(inputFilename string, outputFilename string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFilename,
		"-c:v", "copy",
		"-c:a", "copy",
		"-hls_playlist_type", "vod",
		"-hls_time", "10",
		outputFilename,
	)

	if isLoggingEnabled() {
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func isLoggingEnabled() bool {
	return AppConfiguration.JsonConfigurationProvider{}.ReadRoot().InfrastructureLayerConfiguration.FFMPEGConverter.UseLogging
}

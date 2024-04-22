package FFmpeg

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateSh(t *testing.T) {
	tempDir := t.TempDir()
	scriptPath, err := generateSh(tempDir, "in.mp4")
	if err != nil {
		t.Error(err)
	}

	scriptData, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Error(err)
	}

	if len(scriptData) == 0 {
		t.Error("Empty script data")
	}

	script := string(scriptData)
	if !strings.HasPrefix(script, "#!/bin/bash") {
		t.Error("Invalid script: it does not start with #!/bin/bash")
	}

	if !strings.Contains(script,
		" -c:v libx264"+
			" -c:a aac"+
			" -b:a 128k"+
			" -preset veryfast"+
			" -hls_time 10"+
			" -hls_playlist_type vod"+
			" -hls_playlist_type vod") {
		t.Error("Invalid script")
	}
}

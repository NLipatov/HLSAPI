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
			" -preset ultrafast"+
			" -tune fastdecode"+
			" -crf 35"+
			" -vf scale=1280:-1"+
			" -r 30"+
			" -b:v 2M"+
			" -c:a aac"+
			" -b:a 128k"+
			" -movflags +faststart"+
			" -hls_time 2"+
			" -hls_playlist_type vod"+
			" -hls_key_info_file ") {
		t.Error("Invalid script")
	}
}

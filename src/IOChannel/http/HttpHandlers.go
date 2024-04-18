package httpHandlers

import (
	"fmt"
	"hlsapi/src/Application"
	"hlsapi/src/Infrastructure"
	"hlsapi/src/Infrastructure/FFmpeg"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving: %s\n", r.URL.Query().Get("filename"))
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	filename := r.URL.Query().Get("filename")
	if len(filename) == 0 {
		http.Error(w, "'filename' is mandatory query parameter", http.StatusBadRequest)
		return
	}

	err := Application.GetFileFromDisk(w, filename, Infrastructure.DiskInteractor{})
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
}

func RespondToAHealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Ready"))
	if err != nil {
		panic(err)
	}
}

func CreateM3U8(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	err := r.ParseMultipartForm(300 * 1024 * 1024) // 300 MB max size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadedFiles := r.MultipartForm.File["payload"]
	if len(uploadedFiles) != 1 {
		http.Error(w, "Send 1 file in a form-data with a 'payload' hexKey", http.StatusBadRequest)
		return
	}

	formFile := uploadedFiles[0]
	f, openFileErr := formFile.Open()
	if openFileErr != nil {
		http.Error(w, "Error retrieving formFile from form data", http.StatusBadRequest)
		return
	}

	defer f.Close()

	playlist, err := Application.ConvertVideoToM3U8Playlist(formFile.Filename, f, FFmpeg.Converter{}, Infrastructure.EnvironmentManager{})
	if err != nil {
		http.Error(w, err.Error(), 400)
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(playlist))
	if err != nil {
		panic(err)
	}
}

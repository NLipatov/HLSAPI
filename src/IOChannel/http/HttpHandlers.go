package httpHandlers

import (
	"fmt"
	"hlsapi/src/Application"
	"hlsapi/src/Infrastructure"
	"net/http"
)

func StoreFileOnDisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	err := r.ParseMultipartForm(300 * 1024 * 1024) // 300 MB max size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadedFiles := r.MultipartForm.File["payload"]
	if len(uploadedFiles) == 0 {
		http.Error(w, "Send files in a form-data with a 'payload' key", http.StatusBadRequest)
		return
	}

	for _, formFile := range uploadedFiles {
		func() {
			f, openFileErr := formFile.Open()
			if openFileErr != nil {
				http.Error(w, "Error retrieving formFile from form data", http.StatusBadRequest)
				return
			}

			defer f.Close()

			err = Application.StoreFileOnDisk(formFile.Filename, f, Infrastructure.DiskInteractor{})
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
		}()
	}

	w.WriteHeader(http.StatusCreated)
}

func GetFileFromDisk(w http.ResponseWriter, r *http.Request) {
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
}

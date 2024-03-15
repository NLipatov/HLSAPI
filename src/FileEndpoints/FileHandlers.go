package FileEndpoints

import (
	"bufio"
	"errors"
	"fmt"
	"hlsapi/src/Configuration"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	for _, formfile := range uploadedFiles {
		if CanFileBeStored(formfile.Filename) {
			http.Error(w, "All files should be either .m3u8 or .ts", http.StatusBadRequest)
			return
		}
	}
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

			fileOnDisk, err := os.OpenFile(filepath.Join(Configuration.ReadConfiguration().StorageFolderPath, formFile.Filename), os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, "Error creating formFile in storage", http.StatusInternalServerError)
				return
			}
			defer fileOnDisk.Close()

			_, err = io.Copy(fileOnDisk, f)
			if err != nil {
				http.Error(w, "Error copying formFile data", http.StatusInternalServerError)
				return
			}
		}()
	}

	w.WriteHeader(http.StatusCreated)
}

func CanFileBeStored(filename string) bool {
	isM3U8 := filepath.Ext(filename) != ".m3u8"
	isTs := filepath.Ext(filename) != ".ts"

	return isM3U8 && isTs
}

func GetFileFromDisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	filename := r.URL.Query().Get("filename")
	if len(filename) == 0 {
		http.Error(w, "'filename' is mandatory query parameter", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("%s%c%s", Configuration.ReadConfiguration().StorageFolderPath, os.PathSeparator, filename)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePtr, err := os.Open(path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer filePtr.Close()

	reader := bufio.NewReader(filePtr)
	fInfo, err := filePtr.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = io.CopyN(w, reader, fInfo.Size())

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

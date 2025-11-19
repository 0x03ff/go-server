package json_handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *JsonHandlers) UploadFileHandler(w http.ResponseWriter, r *http.Request) {

	csrf_err := utils.VerifyCSRFtoken(w, r)
	if csrf_err != nil {
		return
	}

	user_id := chi.URLParam(r, "user_id")

	// Validate UUID
	if _, err := uuid.Parse(user_id); err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get user-provided filename (without extension)
	userFilename := filepath.Base(r.FormValue("filename"))
	if userFilename == "" {
		// filename not vaild
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	err = utils.ValidateUserInput("file name", r.FormValue("filename"), 6, 20)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate the filename length (6-20 characters for the base name)
	if len(userFilename) < 6 || len(userFilename) > 20 {
		http.Error(w, "File name must be 6-20 characters", http.StatusBadRequest)
		return
	}

	// Get the file extension from form data (critical fix)
	extension := filepath.Base(r.FormValue("extension"))
	if extension == "" {
		// Fallback to extracting from handler.Filename
		extension = strings.TrimPrefix(filepath.Ext(handler.Filename), ".")
		if extension == "" {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}
	}

	err = utils.ValidateUserInput("file extension", r.FormValue("extension"), 2, 10)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Combine user filename with extension for storage
	filename := userFilename + "." + extension

	// Define the directory path
	dirPath := filepath.Join("assets", "users", user_id, "file")
	// Create the directory if it doesn't exist
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Define the file path
	user_file_path := filepath.Join(dirPath, filename)

	// Clean the path to resolve traversal sequences
	cleanPath := filepath.Clean(user_file_path)

	// 1. Must start with the target directory
	if !strings.HasPrefix(cleanPath, dirPath) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 2. Must not be the directory itself
	if cleanPath == dirPath {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// 3. Verify path separator after directory
	if len(cleanPath) > len(dirPath) && cleanPath[len(dirPath)] != os.PathSeparator {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	filePath := cleanPath

	// Construct the models.File struct from form data
	fileModel := &models.File{
		Title:     filename, // Full name with extension
		UserID:    uuid.MustParse(user_id),
		FilePath:  filePath,
		Extension: extension, // CRITICAL: Set the extension field
	}

	// Store the file information in the database
	fileRepo := repositories.NewFilesRepository(h.dbPool)
	err = fileRepo.Upload(r.Context(), fileModel)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error recording file information: %v", err), http.StatusInternalServerError)
		return
	}

	// Open the destination file for writing
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating file: %v", err), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded successfully"})
}

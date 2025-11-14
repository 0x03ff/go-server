package json_handler

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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

func (h *JsonHandlers) UploadFolderHandler(w http.ResponseWriter, r *http.Request) {

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

	// Parse the multipart form (increase max size for folders)
	err := r.ParseMultipartForm(100 << 20) // 100 MB max folder size
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get folder name from form
	folderName := filepath.Base(r.FormValue("folder_name"))
	if folderName == "" {
		http.Error(w, "Invalid folder name", http.StatusBadRequest)
		return
	}

	err = utils.ValidateUserInput("folder name", folderName, 6, 20)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get encryption method
	encryptMethod := r.FormValue("encrypt_method")
	if encryptMethod == "" {
		encryptMethod = "non-encrypted"
	}

	// Validate encryption method
	validMethods := map[string]bool{
		"non-encrypted": true,
		"aes":           true,
		"rsa-2048":      true,
		"rsa-4096":      true,
	}
	if !validMethods[encryptMethod] {
		http.Error(w, "Invalid encryption method", http.StatusBadRequest)
		return
	}

	// Get all files from the folder
	form := r.MultipartForm
	files := form.File["folder"]
	if len(files) == 0 {
		http.Error(w, "No files in folder", http.StatusBadRequest)
		return
	}

	// Define the directory path
	dirPath := filepath.Join("assets", "users", user_id, "folder", folderName)
	
	// Create the directory if it doesn't exist
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a zip file in memory
	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	// Add all files to zip
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Read file content
		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
			return
		}

		// Add file to zip with relative path
		zipFile, err := zipWriter.Create(fileHeader.Filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating zip entry: %v", err), http.StatusInternalServerError)
			return
		}

		_, err = zipFile.Write(fileContent)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing to zip: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Close the zip writer
	err = zipWriter.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error closing zip: %v", err), http.StatusInternalServerError)
		return
	}

	// Get zip data
	zipData := zipBuffer.Bytes()

	// Encrypt the zip file based on encryption method
	var encryptedData []byte
	var aesKey []byte
	var rsaPrivateKey []byte // Store RSA private key if using RSA encryption

	fmt.Printf("[ENCRYPTION] Method: %s, Original ZIP size: %d bytes\n", encryptMethod, len(zipData))

	switch encryptMethod {
	case "non-encrypted":
		encryptedData = zipData
		fmt.Printf("[ENCRYPTION] No encryption applied (non-encrypted mode)\n")

	case "aes":
		// Generate AES key
		aesKey, err = repositories.GenerateAESKey()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating AES key: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("[ENCRYPTION] Generated AES-256 key: %d bytes\n", len(aesKey))

		// Encrypt with AES
		encryptedData, err = repositories.EncryptAESGCM(zipData, aesKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encrypting with AES: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("[ENCRYPTION] AES encrypted data size: %d bytes (overhead: %d bytes)\n", 
			len(encryptedData), len(encryptedData)-len(zipData))

	case "rsa-2048", "rsa-4096":
		// Generate AES key for hybrid encryption
		aesKey, err = repositories.GenerateAESKey()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating AES key: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("[ENCRYPTION] Generated AES-256 key for hybrid encryption: %d bytes\n", len(aesKey))

		// Encrypt data with AES
		encryptedData, err = repositories.EncryptAESGCM(zipData, aesKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encrypting with AES: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("[ENCRYPTION] AES encrypted data size: %d bytes\n", len(encryptedData))

		// Get RSA key size
		var keySize int
		if encryptMethod == "rsa-2048" {
			keySize = 2048
		} else {
			keySize = 4096
		}
		fmt.Printf("[ENCRYPTION] Using RSA-%d for key encryption\n", keySize)

		// Generate RSA key pair
		rsaPublicKeyPEM, rsaPrivateKeyPEM, err := repositories.GenerateRSAKeyPair(keySize)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating RSA key: %v", err), http.StatusInternalServerError)
			return
		}

		// Parse RSA public key
		block, _ := pem.Decode(rsaPublicKeyPEM)
		if block == nil {
			http.Error(w, "Failed to parse RSA public key", http.StatusInternalServerError)
			return
		}

		pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing public key: %v", err), http.StatusInternalServerError)
			return
		}

		rsaPublicKey := pubKeyInterface.(*rsa.PublicKey)

		// Encrypt AES key with RSA
		encryptedAESKey, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, aesKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encrypting AES key with RSA: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Printf("[ENCRYPTION] RSA encrypted AES key size: %d bytes\n", len(encryptedAESKey))

		// Store encrypted AES key and RSA private key for later decryption
		aesKey = encryptedAESKey
		rsaPrivateKey = rsaPrivateKeyPEM
	}

	fmt.Printf("[ENCRYPTION] Final encrypted data size: %d bytes, Method: %s\n", len(encryptedData), encryptMethod)

	// Define the file path for the encrypted zip
	zipFileName := folderName + ".zip"
	zipFilePath := filepath.Join(dirPath, zipFileName)

	// Clean the path to resolve traversal sequences
	cleanPath := filepath.Clean(zipFilePath)

	// Security checks
	if !strings.HasPrefix(cleanPath, dirPath) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Write encrypted data to file
	err = os.WriteFile(cleanPath, encryptedData, 0644)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving encrypted folder: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("[ENCRYPTION] Saved encrypted file to: %s\n", cleanPath)
	
	// Verify what was written
	fileInfo, _ := os.Stat(cleanPath)
	fmt.Printf("[ENCRYPTION] File on disk size: %d bytes\n", fileInfo.Size())
	
	// Check first few bytes to verify encryption
	savedData, _ := os.ReadFile(cleanPath)
	if len(savedData) >= 4 {
		fmt.Printf("[ENCRYPTION] First 4 bytes (hex): %02x %02x %02x %02x\n", 
			savedData[0], savedData[1], savedData[2], savedData[3])
		if savedData[0] == 0x50 && savedData[1] == 0x4B {
			fmt.Printf("[WARNING] File starts with ZIP magic bytes (PK) - NOT ENCRYPTED!\n")
		} else {
			fmt.Printf("[SUCCESS] File does NOT start with ZIP magic bytes - ENCRYPTED!\n")
		}
	}

	// Construct the models.Folder struct
	folderModel := &models.Folder{
		Title:      folderName,
		UserID:     uuid.MustParse(user_id),
		FilePath:   cleanPath,
		Secret:     aesKey,        // Store AES key (or encrypted AES key for RSA)
		PrivateKey: rsaPrivateKey, // Store RSA private key (only for RSA encryption)
		Encrypt:    encryptMethod,
		Extension:  "zip",
	}

	// Store the folder information in the database
	folderRepo := repositories.NewFoldersRepository(h.dbPool)
	err = folderRepo.Upload(r.Context(), folderModel)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error recording folder information: %v", err), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Folder uploaded successfully"})
}
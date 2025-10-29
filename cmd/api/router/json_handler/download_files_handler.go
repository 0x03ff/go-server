package json_handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (h *JsonHandlers) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
    userIDParam := chi.URLParam(r, "user_id")
    fileIDParam := chi.URLParam(r, "file_id")

    if userIDParam == "" || fileIDParam == "" {
        http.Error(w, "User ID or File ID not found in URL path", http.StatusBadRequest)
        return
    }

    // Get the token from the cookie
    cookie, err := r.Cookie("token")
    if err != nil {
        http.Error(w, "Token not found. Please log in again.", http.StatusUnauthorized)
        return
    }

    token := cookie.Value
    if token == "" {
        http.Error(w, "Invalid token format", http.StatusUnauthorized)
        return
    }

    systemRepo := repositories.NewKeysRepository(h.dbPool)
    // Verify the token
    tokenObj, err := utils.VerifyToken(token, systemRepo)
    if err != nil {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    // Extract claims from the token
    claims := tokenObj.Claims.(jwt.MapClaims)
    userIDClaim, ok := claims["user_id"].(string)
    if !ok {
        http.Error(w, "Invalid token claims", http.StatusUnauthorized)
        return
    }

    // Check if the user_id in the URL matches the user_id in the token
    if userIDClaim != userIDParam {
        http.Error(w, "User ID mismatch", http.StatusUnauthorized)
        return
    }

    // Parse user_id to UUID
    parsedUserID, err := uuid.Parse(userIDParam)
    if err != nil {
        http.Error(w, "Invalid user ID format", http.StatusBadRequest)
        return
    }

    // Parse file_id to integer
    fileID, err := strconv.Atoi(fileIDParam)
    if err != nil {
        http.Error(w, "Invalid file ID format", http.StatusBadRequest)
        return
    }

    // Get the file from the database
    filesRepo := repositories.NewFilesRepository(h.dbPool)
    file := &models.File{}
    err = filesRepo.GetFileById(r.Context(), file, fileID)
    if err != nil {
        http.Error(w, "Failed to retrieve file information", http.StatusInternalServerError)
        return
    }

    // Verify the file belongs to the user
    if file.UserID != parsedUserID {
        http.Error(w, "You don't have permission to download this file", http.StatusForbidden)
        return
    }

    // CRITICAL: Set correct headers for binary download
    // DO NOT use base64 encoding - send raw binary
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.enc\"", file.Title))
    w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
    w.Header().Set("X-Content-Type-Options", "nosniff")

    // Open the file
    fileData, err := os.Open(file.FilePath)
    if err != nil {
        http.Error(w, "Failed to open file", http.StatusInternalServerError)
        return
    }
    defer fileData.Close()

    // Stream the file directly to response - NO ENCODING!
    if _, err := io.Copy(w, fileData); err != nil {
        http.Error(w, "Failed to send file", http.StatusInternalServerError)
        return
    }
}

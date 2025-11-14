package json_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (h *JsonHandlers) ObtainFolderHandler(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")

	if user_id == "" {
		http.Error(w, "User ID not found in URL path", http.StatusBadRequest)
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
	userIdClaim, ok := claims["user_id"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Check if the user_id in the URL matches the user_id in the token
	if userIdClaim != user_id {
		http.Error(w, "User ID mismatch", http.StatusUnauthorized)
		return
	}
	
	// Parse user_id to UUID
	parsedUserID, err := uuid.Parse(user_id)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	
	// Get index from query parameter (default to 0)
	indexStr := r.URL.Query().Get("index")
	index := 0
	if indexStr != "" {
		index, err = strconv.Atoi(indexStr)
		if err != nil || index < 0 {
			http.Error(w, "Index must be a non-negative integer", http.StatusBadRequest)
			return
		}
	}
	
	// Call the repository to search for folders
	foldersRepo := repositories.NewFoldersRepository(h.dbPool)
	folders, err := foldersRepo.Search(r.Context(), parsedUserID, index)
	if err != nil {
		http.Error(w, "Failed to retrieve folders", http.StatusInternalServerError)
		return
	}
	
	// Check if there are more folders after this page
	hasNext := len(folders) == 10
	
	// Check if there are folders before this page
	hasPrev := index > 0
	
	// Calculate current page number (1-based)
	page := (index / 10) + 1

	// Prepare response using the models.FoldersResponse structure
	response := models.FoldersResponse{
		Folders: folders,
		Index:   index,
		Page:    page,
		HasPrev: hasPrev,
		HasNext: hasNext,
	}
	
	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

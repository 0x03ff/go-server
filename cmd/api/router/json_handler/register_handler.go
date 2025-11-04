package json_handler

import (
	"encoding/json"
	"net/http"

	"github.com/0x03ff/golang/internal/store/models"
	"github.com/0x03ff/golang/internal/store/repositories"
	"github.com/0x03ff/golang/utils"
)

func (h *JsonHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	
	csrf_err := utils.VerifyCSRFtoken(w , r )
	if csrf_err != nil{
		return
	}
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	// Validate user data
	if user.Username == "" || user.Password == "" {
		utils.SendError(w, http.StatusBadRequest, "Invalid user data")
		return
	}

	err = utils.ValidateUserInput("User ID", user.Username, 6, 20)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = utils.ValidateUserInput("Password", user.Password, 8, 20)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}	
	err = utils.ValidateUserInput("Recover key", user.Recover, 6, 20)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	

	userRepo := repositories.NewUsersRepository(h.dbPool)
	err = userRepo.Create(r.Context(), &user)

	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Registration failed. Please try again.")
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
	})
}

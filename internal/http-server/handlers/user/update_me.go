package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (i *Implementation) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userIDStr, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		PhotoURL  string `json:"photo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Username == "" && req.FirstName == "" && req.PhotoURL == "" {
		response.WriteError(w, http.StatusBadRequest, "At least one field (username, first_name or photo_url) must be provided")
		return
	}

	user := &models.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		PhotoURL:  req.PhotoURL,
	}

	if err := i.s.Update(r.Context(), userID, user); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response.WriteSuccess(w, http.StatusOK, "User updated successfully")
}


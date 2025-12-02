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

func (i *Implementation) Update(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	path = strings.TrimSuffix(path, "/")

	if path == "" || path == r.URL.Path {
		response.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if !middleware.CheckUserOwnership(r, id) {
		response.WriteError(w, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var req struct {
		Username string `json:"username"`
		PhotoURL string `json:"photo_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Проверяем, что хотя бы одно поле передано
	if req.Username == "" && req.PhotoURL == "" {
		response.WriteError(w, http.StatusBadRequest, "At least one field (username or photo_url) must be provided")
		return
	}

	user := &models.User{
		Username: req.Username,
		PhotoURL: req.PhotoURL,
	}

	if err := i.s.Update(r.Context(), id, user); err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	response.WriteSuccess(w, http.StatusOK, "User updated successfully")
}

package settings

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) UpdateTimezoneMe(w http.ResponseWriter, r *http.Request) {
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
		Timezone string `json:"timezone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(req.Timezone) == "" {
		response.WriteError(w, http.StatusBadRequest, "timezone cannot be empty")
		return
	}

	updatedSettings, err := i.s.UpdateTimezone(r.Context(), userID, req.Timezone)
	if err != nil {
		if strings.Contains(err.Error(), "timezone validation failed") || strings.Contains(err.Error(), "invalid timezone") {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update timezone")
		return
	}

	response.WriteJSON(w, http.StatusOK, updatedSettings)
}







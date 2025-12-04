package settings

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
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
		DoNotDisturb *bool     `json:"do_not_disturb,omitempty"`
		NotifyTimes  *[]string `json:"notify_times,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.DoNotDisturb == nil && req.NotifyTimes == nil {
		response.WriteError(w, http.StatusBadRequest, "At least one field (do_not_disturb or notify_times) must be provided")
		return
	}

	updatedSettings, err := i.s.UpdateSettings(r.Context(), userID, req.DoNotDisturb, req.NotifyTimes)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	response.WriteJSON(w, http.StatusOK, updatedSettings)
}


package settings

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) Update(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	path = strings.TrimSuffix(path, "/settings")

	if path == "" || path == r.URL.Path {
		response.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	path = strings.TrimSuffix(path, "/")

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if !middleware.CheckUserOwnership(r, id) {
		response.WriteError(w, http.StatusForbidden, "You can only update your own settings")
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

	updatedSettings, err := i.s.UpdateSettings(r.Context(), id, req.DoNotDisturb, req.NotifyTimes)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	response.WriteJSON(w, http.StatusOK, updatedSettings)
}

package sprint

import (
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) GetAll(w http.ResponseWriter, r *http.Request) {
	var isActive *bool
	if activeStr := r.URL.Query().Get("is_active"); activeStr != "" {
		active, err := strconv.ParseBool(activeStr)
		if err == nil {
			isActive = &active
		}
	}

	sprints, err := i.sprintService.GetAll(r.Context(), isActive)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, sprints)
}


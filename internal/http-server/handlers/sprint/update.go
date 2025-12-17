package sprint

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (i *Implementation) Update(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sprints/")
	sprintID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid sprint ID")
		return
	}

	var req models.CreateSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	sprint, err := i.sprintService.Update(r.Context(), sprintID, &req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, sprint)
}


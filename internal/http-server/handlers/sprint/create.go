package sprint

import (
	"encoding/json"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSprintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	sprint, err := i.sprintService.Create(r.Context(), &req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, sprint)
}


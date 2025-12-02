package habit

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
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
		Title        string           `json:"title"`
		Type         models.HabitType `json:"type"`
		Unit         string           `json:"unit"`
		Value        int              `json:"value"`
		IsActive     *bool            `json:"is_active,omitempty"`
		IsBeneficial *bool            `json:"is_beneficial,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	habit := &models.Habit{
		UserID: userID,
		Title:  req.Title,
		Type:   req.Type,
		Unit:   req.Unit,
		Value:  req.Value,
	}

	if req.IsActive != nil {
		habit.IsActive = *req.IsActive
	} else {
		habit.IsActive = true
	}

	if req.IsBeneficial != nil {
		habit.IsBeneficial = *req.IsBeneficial
	} else {
		habit.IsBeneficial = true // По умолчанию полезная привычка
	}

	createdHabit, err := i.s.Create(r.Context(), habit)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, createdHabit)
}

package habit

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

	path := strings.TrimPrefix(r.URL.Path, "/habits/")
	if path == "" || path == r.URL.Path {
		response.WriteError(w, http.StatusBadRequest, "Habit ID is required")
		return
	}

	path = strings.TrimSuffix(path, "/")

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid habit ID")
		return
	}

	currentHabit, err := i.s.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.WriteError(w, http.StatusNotFound, "Habit not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to get habit")
		return
	}

	if currentHabit.UserID != userID {
		response.WriteError(w, http.StatusForbidden, "You can only update your own habits")
		return
	}

	var req struct {
		Title        *string           `json:"title,omitempty"`
		Type         *models.HabitType `json:"type,omitempty"`
		Unit         *string           `json:"unit,omitempty"`
		Value        *int              `json:"value,omitempty"`
		IsActive     *bool             `json:"is_active,omitempty"`
		IsDone       *bool             `json:"is_done,omitempty"`
		IsBeneficial *bool             `json:"is_beneficial,omitempty"`
		Series       *int              `json:"series,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Title != nil {
		currentHabit.Title = *req.Title
	}
	if req.Type != nil {
		currentHabit.Type = *req.Type
	}
	if req.Unit != nil {
		currentHabit.Unit = *req.Unit
	}
	if req.Value != nil {
		currentHabit.Value = *req.Value
	}
	if req.IsActive != nil {
		currentHabit.IsActive = *req.IsActive
	}
	if req.IsDone != nil {
		currentHabit.IsDone = *req.IsDone
	}
	if req.IsBeneficial != nil {
		currentHabit.IsBeneficial = *req.IsBeneficial
	}
	if req.Series != nil {
		currentHabit.Series = *req.Series
	}

	if err := i.s.Update(r.Context(), currentHabit); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update habit")
		return
	}

	response.WriteJSON(w, http.StatusOK, currentHabit)
}

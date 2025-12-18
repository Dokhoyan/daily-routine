package habit

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/logger"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("habit create handler: started")

	userIDStr, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.Errorf("habit create handler: failed to get user id from context: %v", err)
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	logger.Infof("habit create handler: user id from context: %s", userIDStr)

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Errorf("habit create handler: failed to parse user id: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	logger.Infof("habit create handler: parsed user id: %d", userID)

	var req struct {
		Title        string             `json:"title"`
		Format       models.HabitFormat `json:"format"`
		Unit         string             `json:"unit"`
		Value        int                `json:"value"`
		CurrentValue *int               `json:"current_value,omitempty"`
		IsActive     *bool              `json:"is_active,omitempty"`
		Type         *models.HabitType  `json:"type,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf("habit create handler: failed to decode json: %v", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	logger.Infof("habit create handler: request decoded: title=%s, format=%s, unit=%s, value=%d", req.Title, req.Format, req.Unit, req.Value)

	habit := &models.Habit{
		UserID: userID,
		Title:  req.Title,
		Format: req.Format,
		Unit:   req.Unit,
		Value:  req.Value,
	}

	if req.CurrentValue != nil {
		habit.CurrentValue = *req.CurrentValue
	} else {
		habit.CurrentValue = 0
	}

	if req.IsActive != nil {
		habit.IsActive = *req.IsActive
	} else {
		habit.IsActive = true
	}

	if req.Type != nil {
		habit.Type = *req.Type
	} else {
		habit.Type = models.HabitTypeBeneficial
	}

	if habit.Type == models.HabitTypeHarmful {
		habit.IsDone = true
	} else {
		habit.IsDone = false
	}

	habit.Series = 0

	logger.Infof("habit create handler: calling service with habit: userID=%d, title=%s, format=%s, type=%s, value=%d, currentValue=%d, isActive=%v, isDone=%v",
		habit.UserID, habit.Title, habit.Format, habit.Type, habit.Value, habit.CurrentValue, habit.IsActive, habit.IsDone)

	createdHabit, err := i.s.Create(r.Context(), habit)
	if err != nil {
		logger.Errorf("habit create handler: service error: %v", err)
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof("habit create handler: habit created successfully with id=%d", createdHabit.ID)
	response.WriteJSON(w, http.StatusCreated, createdHabit)
}

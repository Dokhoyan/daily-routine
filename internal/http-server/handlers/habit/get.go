package habit

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) {
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

	habit, err := i.s.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.WriteError(w, http.StatusNotFound, "Habit not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to get habit")
		return
	}

	if habit.UserID != userID {
		response.WriteError(w, http.StatusForbidden, "You can only access your own habits")
		return
	}

	response.WriteJSON(w, http.StatusOK, habit)
}

func (i *Implementation) GetAll(w http.ResponseWriter, r *http.Request) {
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

	var habitType *string
	var isActive *bool

	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		if typeParam != "beneficial" && typeParam != "harmful" {
			response.WriteError(w, http.StatusBadRequest, "Invalid type parameter: must be 'beneficial' or 'harmful'")
			return
		}
		habitType = &typeParam
	}

	if isActiveParam := r.URL.Query().Get("is_active"); isActiveParam != "" {
		parsed, err := strconv.ParseBool(isActiveParam)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "Invalid is_active parameter: must be 'true' or 'false'")
			return
		}
		isActive = &parsed
	}

	habits, err := i.s.GetByUserID(r.Context(), userID, habitType, isActive)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get habits")
		return
	}

	response.WriteJSON(w, http.StatusOK, habits)
}

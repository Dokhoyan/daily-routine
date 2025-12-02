package habit

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) Delete(w http.ResponseWriter, r *http.Request) {
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
		response.WriteError(w, http.StatusInternalServerError, "Failed to get habit for ownership check")
		return
	}

	if habit.UserID != userID {
		response.WriteError(w, http.StatusForbidden, "You can only delete your own habits")
		return
	}

	if err := i.s.Delete(r.Context(), id); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete habit")
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Habit deleted successfully")
}

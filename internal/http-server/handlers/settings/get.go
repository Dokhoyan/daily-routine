package settings

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/middleware"
	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) {
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
		response.WriteError(w, http.StatusForbidden, "You can only access your own settings")
		return
	}

	settings, err := i.s.GetByUserID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get settings")
		return
	}

	response.WriteJSON(w, http.StatusOK, settings)
}

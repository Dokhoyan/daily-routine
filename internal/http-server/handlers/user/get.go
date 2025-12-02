package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) Get(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	if path == "" || path == r.URL.Path {
		response.WriteError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	path = strings.TrimSuffix(path, "/")

	if strings.Contains(path, "/") {
		response.WriteError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := i.s.GetByID(r.Context(), id)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

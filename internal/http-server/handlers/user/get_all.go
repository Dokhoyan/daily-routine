package user

import (
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := i.s.GetAll(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	response.WriteJSON(w, http.StatusOK, users)
}

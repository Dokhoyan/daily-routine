package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) LoginOrRegister(w http.ResponseWriter, r *http.Request) {
	var telegramData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&telegramData); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	authResponse, err := i.authService.AuthenticateOrRegister(r.Context(), telegramData, r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, authResponse)
}

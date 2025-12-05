package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
	"github.com/Dokhoyan/daily-routine/internal/models"
)

func (i *Implementation) GetAccessToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.RefreshToken == "" {
		response.WriteError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	accessToken, err := i.authService.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	tokenPair := &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	}

	response.WriteJSON(w, http.StatusOK, tokenPair)
}

func (i *Implementation) GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.RefreshToken == "" {
		response.WriteError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	tokenPair, err := i.authService.RefreshTokenPair(r.Context(), req.RefreshToken, r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, tokenPair)
}

package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) TestToken(enableTestMode bool, defaultTestUserID int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !enableTestMode {
			response.WriteError(w, http.StatusForbidden, "Test mode is disabled")
			return
		}

		var req struct {
			UserID int64 `json:"user_id,omitempty"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		userID := defaultTestUserID
		if req.UserID != 0 {
			userID = req.UserID
		}

		testResponse, err := i.authService.GenerateTestToken(r.Context(), userID, r)
		if err != nil {
			response.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJSON(w, http.StatusOK, testResponse)
	}
}

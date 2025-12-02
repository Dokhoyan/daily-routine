package auth

import (
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/service"
)

type AuthService = service.AuthService

type Implementation struct {
	authService AuthService
}

func NewImplementation(authService AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login/telegram", i.LoginOrRegister)
	mux.HandleFunc("/auth/getaccesstoken", i.GetAccessToken)
	mux.HandleFunc("/auth/getrefreshtoken", i.GetRefreshToken)
	mux.HandleFunc("/auth/revoke", i.RevokeToken)
}

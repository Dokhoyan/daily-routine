package user

import (
	"net/http"

	"github.com/Dokhoyan/daily-routine/internal/service"
)

type UserService = service.UserService

type Implementation struct {
	s UserService
}

func NewImplementation(s UserService) *Implementation {
	return &Implementation{
		s: s,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", i.GetAll)
	mux.HandleFunc("/user/me", i.handleMeRoutes)
}

func (i *Implementation) handleMeRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		i.GetMe(w, r)
	} else if r.Method == http.MethodPut || r.Method == http.MethodPatch {
		i.UpdateMe(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

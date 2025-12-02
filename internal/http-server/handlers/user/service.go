package user

import (
	"net/http"
	"strings"

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
	mux.HandleFunc("/users/", i.handleUserRoutes)
}

func (i *Implementation) HandleUserRoutes(w http.ResponseWriter, r *http.Request) {
	i.handleUserRoutes(w, r)
}

func (i *Implementation) handleUserRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	path = strings.TrimSuffix(path, "/")

	if path == "" || path == r.URL.Path {
		i.Get(w, r)
		return
	}

	if r.Method == http.MethodPut || r.Method == http.MethodPatch {
		i.Update(w, r)
		return
	}

	i.Get(w, r)
}

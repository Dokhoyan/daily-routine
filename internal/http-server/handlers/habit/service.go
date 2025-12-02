package habit

import (
	"net/http"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/service"
)

type HabitService = service.HabitService

type Implementation struct {
	s HabitService
}

func NewImplementation(s HabitService) *Implementation {
	return &Implementation{
		s: s,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/habits", i.handleHabitsRoutes)
	mux.HandleFunc("/habits/", i.handleHabitRoutes)
}

func (i *Implementation) handleHabitsRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		i.Create(w, r)
	} else if r.Method == http.MethodGet {
		i.GetAll(w, r)
	}
}

func (i *Implementation) handleHabitRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/habits/")

	if path == "" || path == r.URL.Path {
		return
	}

	path = strings.TrimSuffix(path, "/")

	switch r.Method {
	case http.MethodGet:
		i.Get(w, r)
	case http.MethodPut, http.MethodPatch:
		i.Update(w, r)
	case http.MethodDelete:
		i.Delete(w, r)
	}
}

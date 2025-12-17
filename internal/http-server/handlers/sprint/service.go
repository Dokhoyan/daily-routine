package sprint

import (
	"net/http"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/service"
)

type Implementation struct {
	sprintService service.SprintService
	userService   service.UserService
}

func NewImplementation(sprintService service.SprintService, userService service.UserService) *Implementation {
	return &Implementation{
		sprintService: sprintService,
		userService:   userService,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/sprints", i.handleSprintsRoutes)
	mux.HandleFunc("/sprints/", i.handleSprintRoutes)
	mux.HandleFunc("/sprints/progress", i.handleProgressRoute)
}

func (i *Implementation) handleSprintsRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		i.GetAll(w, r)
	} else if r.Method == http.MethodPost {
		i.Create(w, r)
	}
}

func (i *Implementation) handleSprintRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sprints/")

	if path == "" || path == r.URL.Path {
		return
	}

	path = strings.TrimSuffix(path, "/")

	switch r.Method {
	case http.MethodGet:
		// Можно добавить GetByID если нужно
		i.GetAll(w, r)
	case http.MethodPut, http.MethodPatch:
		i.Update(w, r)
	case http.MethodDelete:
		i.Delete(w, r)
	}
}

func (i *Implementation) handleProgressRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		i.GetProgress(w, r)
	}
}


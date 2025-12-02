package settings

import (
	"net/http"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/service"
)

type SettingsService = service.SettingsService

type Implementation struct {
	s SettingsService
}

func NewImplementation(s SettingsService) *Implementation {
	return &Implementation{
		s: s,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/", i.handleSettingsRoutes)
}

func (i *Implementation) HandleSettingsRoutes(w http.ResponseWriter, r *http.Request) {
	i.handleSettingsRoutes(w, r)
}

func (i *Implementation) handleSettingsRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/users/")

	if strings.HasSuffix(path, "/settings/timezone") {
		if r.Method == http.MethodPut {
			i.UpdateTimezone(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if strings.HasSuffix(path, "/settings") {
		if r.Method == http.MethodGet {
			i.Get(w, r)
		} else if r.Method == http.MethodPatch {
			i.Update(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
}

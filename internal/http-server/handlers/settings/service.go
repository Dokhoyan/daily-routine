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
	mux.HandleFunc("/user/me/settings", i.handleMeSettingsRoutes)
	mux.HandleFunc("/user/me/settings/", i.handleMeSettingsRoutes)
}

func (i *Implementation) handleMeSettingsRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/user/me/settings")
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimPrefix(path, "/")

	if path == "timezone" {
		if r.Method == http.MethodPut {
			i.UpdateTimezoneMe(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if path == "" {
		if r.Method == http.MethodGet {
			i.GetMe(w, r)
		} else if r.Method == http.MethodPatch || r.Method == http.MethodPut {
			i.UpdateMe(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

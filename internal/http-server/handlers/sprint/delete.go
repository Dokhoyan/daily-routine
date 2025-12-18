package sprint

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Dokhoyan/daily-routine/internal/http-server/response"
)

func (i *Implementation) Delete(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sprints/")
	sprintID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid sprint ID")
		return
	}

	if err := i.sprintService.Delete(r.Context(), sprintID); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Sprint deleted successfully")
}



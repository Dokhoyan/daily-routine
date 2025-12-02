package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CheckUserOwnership(r *http.Request, resourceUserID int64) bool {
	userIDStr, err := GetUserIDFromContext(r.Context())
	if err != nil {
		return false
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return false
	}

	return userID == resourceUserID
}

func GetUserIDFromPath(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("id not found in path")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid id format: %w", err)
	}

	return id, nil
}

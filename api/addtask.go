package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"final_project/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %s", err.Error()))
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "title not set")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}

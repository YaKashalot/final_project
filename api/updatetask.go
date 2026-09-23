package api

import (
	"encoding/json"
	"net/http"

	"final_project/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, http.StatusBadRequest, "id is not set")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "title is not set")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{})
}

package api

import (
	"encoding/json"
	"net/http"

	"final_project/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, "invalid JSON: "+err.Error())
		return
	}

	if task.ID == "" {
		writeError(w, "id not set")
		return
	}

	if task.Title == "" {
		writeError(w, "title not set")
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, map[string]any{})
}

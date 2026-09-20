package api

import (
	"net/http"

	"final_project/db"
)

const tasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	var tasks []*db.Task
	var err error

	switch {
	case search == "":
		tasks, err = db.Tasks(tasksLimit)
	default:
		if date, ok := db.ParseSearchDate(search); ok {
			tasks, err = db.TasksByDate(date, tasksLimit)
		} else {
			tasks, err = db.TasksBySearch(search, tasksLimit)
		}
	}

	if err != nil {
		writeError(w, err.Error())
		return
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}

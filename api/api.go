package api

import (
	"net/http"

	"final_project/config"
)

var appConfig config.Config

func Init(cfg config.Config) {
	appConfig = cfg

	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signInHandler)

	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
}

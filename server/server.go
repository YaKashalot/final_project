package server

import (
	"log"
	"net/http"
	"os"

	"final_project/api"
)

const defaultPort = "7540"
const webDir = "./web"

func CreateServer() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Starting server on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

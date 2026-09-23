package server

import (
	"final_project/api"
	"final_project/config"
	"log"
	"net/http"
)

const defaultPort = "7540"
const webDir = "./web"

func CreateServer(cfg config.Config) {
	api.Init(cfg)
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Starting server on port %s...", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

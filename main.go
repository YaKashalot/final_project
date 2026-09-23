package main

import (
	"final_project/config"
	"log"
	"os"

	"final_project/db"
	"final_project/server"
)

func main() {
	cfg := config.Load()

	if err := db.Init(cfg.DBFile); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	server.CreateServer(cfg)
}

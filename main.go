package main

import (
	"log"
	"os"

	"final_project/db"
	"final_project/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	err := db.Init(dbFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	server.CreateServer()
}

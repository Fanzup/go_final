package main

import (
	"final/v1/pkg/db"
	"final/v1/pkg/server"
	"log"
	"os"
)

const defaultFileDB = "scheduler.db"

func main() {

	if err := db.InitDB(findDB()); err != nil {
		log.Fatal("DB is not init! ", err)
	}
	log.Println("DB is init!")

	server.ServerRun()
}

func findDB() string {
	if checkDB := os.Getenv("TODO_DBFILE"); checkDB != "" {
		return checkDB
	}
	return defaultFileDB
}

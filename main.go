package main

import (
	"log"
	"os"

	"final/v1/pkg/db"
	"final/v1/pkg/server"
)

const defaultFileDB = "scheduler.db"

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// поднимаем БД
	if err := db.InitDB(findDB()); err != nil {
		return err
	}
	log.Println("DB is init!")

	defer db.Close()
	return server.ServerRun()
}

func findDB() string {
	if checkDB := os.Getenv("TODO_DBFILE"); checkDB != "" {
		return checkDB
	}
	return defaultFileDB
}

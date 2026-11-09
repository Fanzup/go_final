package server

import (
	"final/v1/pkg/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultPort = 7540

func ServerRun() {
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(defaultPort)
	}

	addr := ":" + port
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Server is running - port=%s", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}

}

package server

import (
	"final/v1/pkg/api"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultPort = 7540

func ServerRun() {
	mux := http.NewServeMux()
	api.Init(mux)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(defaultPort)
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		log.Fatalf("invalid TODO_PORT: %q", port)
	}
	addr := ":" + strconv.Itoa(portNumber)
	webDir := "./web"

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Server is running - port=%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}

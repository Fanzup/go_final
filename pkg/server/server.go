package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"final/v1/pkg/api"
)

const defaultPort = 7540

func ServerRun() error {
	mux := http.NewServeMux()
	if err := api.Init(mux); err != nil {
		return err
	}

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = strconv.Itoa(defaultPort)
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil || portNumber < 1 || portNumber > 65535 {
		return fmt.Errorf("invalid TODO_PORT: %q", port)
	}
	addr := ":" + strconv.Itoa(portNumber)
	webDir := "./web"

	// отдаём файлы фронтенда
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.ListenAndServe() }()
	log.Printf("Server is running - port=%s", port)
	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		// даём активным запросам закончиться перед закрытием БД
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return err
		}
		return nil
	}
}

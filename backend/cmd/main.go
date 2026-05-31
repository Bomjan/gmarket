package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bomjan/gmarket/backend/internal/config"
	"github.com/Bomjan/gmarket/backend/internal/http/handlers/student"
)

func main() {
	//load config

	cfg := config.MustLoadConfig()
	//database setup
	// set up router
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())
	// set up server

	server := http.Server{
		Addr:    cfg.Addr, // Addr:    "localhost:8080",
		Handler: router,
	}

	slog.Info("The server has started", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("failed to start the server")
		}
	}()
	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")
}

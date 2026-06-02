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
	"github.com/Bomjan/gmarket/backend/internal/http/handlers/products"
	"github.com/Bomjan/gmarket/backend/internal/http/handlers/student"
	"github.com/Bomjan/gmarket/backend/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoadConfig()

	//database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("database initialized successfully")

	// set up router
	router := http.NewServeMux()
	// Student api
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetAllStudents(storage))

	// Product api
	router.HandleFunc("POST /api/products", products.New(storage))
	router.HandleFunc("GET /api/products/{id}", products.GetById(storage))
	router.HandleFunc("GET /api/products", products.GetProducts(storage))
	router.HandleFunc("DELETE /api/products/{id}", products.DeleteProductById(storage))
	router.HandleFunc("PUT /api/products/{id}", products.UpdateProductById(storage))

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

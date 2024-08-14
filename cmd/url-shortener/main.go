package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http-server/handlers/url/save"
	"awesomeProject/internal/lib/logger/setup"
	"awesomeProject/internal/lib/router"
	"awesomeProject/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	collectionName = "url-shortener"
)

func main() {
	// CONTEXT
	contextOur := context.TODO()
	// CONFIG
	cfg := config.MustLoad()

	// LOG
	logger := setup.SetupLogger(cfg.Env)
	if logger == nil {
		fmt.Println("Failed to load config")
	}
	logger.Info("Starting logger", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// MONGODB
	mongoStorage, err := storage.ConnectToDB(collectionName, cfg, logger, contextOur)
	if err != nil {
		logger.Error("Can't connect to MongoDB", slog.String("error", err.Error()))
		return
	}
	// ROUTER
	NewRouter := router.SetupRouter(logger)

	NewRouter.Post("/url", save.SaveUrlHandler(logger, mongoStorage))
	// SERVER
	logger.Info("Starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      NewRouter,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}

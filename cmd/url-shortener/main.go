package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/lib/logger/setup"
	"awesomeProject/internal/lib/router"
	"awesomeProject/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// CONTEXT
	ctx := context.TODO()
	// CONFIG
	cfg := config.MustLoad()

	// LOG
	logger := setup.LoggerSetup(cfg.Env)
	if logger == nil {
		fmt.Println("Failed to load config")
	}
	logger.Info("Starting logger", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// Run migrations
	if err := storage.RunMigrations(cfg.DBName, cfg); err != nil {
		logger.Error("Failed to run migrations", slog.String("error", err.Error()))
	}

	// MONGODB
	mongoStorage, err := storage.ConnectToDB(cfg.CollectionName, cfg.DBName, cfg, logger, ctx)
	if err != nil {
		logger.Error("Can't connect to MongoDB", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// ROUTER
	newRouter := router.SetupRouter(logger, mongoStorage, cfg)
	// SERVER
	logger.Info("Starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      newRouter,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server", slog.String("error", err.Error()))
	}

	logger.Info("Server stopped")
}

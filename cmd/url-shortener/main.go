package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/storage"
	"fmt"
	"log/slog"
	"os"
)

const (
	envLocal       = "local"
	envDev         = "dev"
	envProd        = "prod"
	CollectionName = "url-shortener"
)

func main() {
	const op = "cmd/url-shortener/main"
	//CONFIG
	cfg := config.MustLoad()

	//LOG

	log := setupLogger(cfg.Env)
	if log == nil {
		fmt.Println("Failed to load config")
	}
	log.Info("Starting logger", slog.String("env", cfg.Env))
	log.Debug("Starting debug", slog.String("env", cfg.Env))

	// MONGODB

	MongoStorage, err := storage.ConnectionToDB(CollectionName, cfg)
	if err != nil {
		log.Info("Error connecting to Mongo", slog.String("error", err.Error()))
	}
	_ = MongoStorage
	// TODO: init router: chi "chi render"
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

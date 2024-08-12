package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http-server/handlers/url/save"
	mwLogger "awesomeProject/internal/http-server/middleware/logger"
	"awesomeProject/internal/lib/logger/handlers/slogpretty"
	"awesomeProject/internal/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal       = "local"
	envDev         = "dev"
	envProd        = "prod"
	collectionName = "url-shortener"
)

func main() {
	//CONFIG
	cfg := config.MustLoad()

	//LOG
	logger := setupLogger(cfg.Env)
	if logger == nil {
		fmt.Println("Failed to load config")
	}
	logger.Info("Starting logger", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")

	// MONGODB
	mongoStorage, err := storage.ConnectingToDB(collectionName, cfg, logger)
	if err != nil {
		logger.Info("Error connecting to Mongo", slog.String("error", err.Error()))
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID) // мидлвейр с присваиванием АЙДИ каждому запросу, что может помочь при отслеживании
	// действий для конкретного запроса (например при ошибке, чтоб понять, где ошибка)
	router.Use(middleware.Logger) //Свой встроенный логгер в мидлвейр от CHI, будет красиво логгировать каждый запрос клиента, со временем на обработку и т.д.
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer) // в случае паники восстанавливаем, чтоб ничего не падало из-за 1 запроса
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(logger, mongoStorage))

	logger.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("Failed to start server")
	}

	logger.Error("server stopped")
	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog() //Почему то блять не работает подсветка в консоли суки
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

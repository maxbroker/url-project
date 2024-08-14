package router

import (
	mwLogger "awesomeProject/internal/http-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func SetupRouter(logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // мидлвейр с присваиванием АЙДИ каждому запросу, что может помочь при отслеживании
	// действий для конкретного запроса (например при ошибке, чтоб понять, где ошибка)
	router.Use(middleware.Logger) //Свой встроенный логгер в мидлвейр от CHI, будет красиво логгировать каждый запрос клиента, со временем на обработку и т.д.
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer) // в случае паники восстанавливаем, чтоб ничего не падало из-за 1 запроса
	router.Use(middleware.URLFormat)
	return router
}

package router

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http-server/handlers/redirect"
	deleteReq "awesomeProject/internal/http-server/handlers/url/delete"
	"awesomeProject/internal/http-server/handlers/url/save"
	mwLogger "awesomeProject/internal/http-server/middleware/logger"
	"awesomeProject/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func SetupRouter(logger *slog.Logger, storage *storage.Storage, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // мидлвейр с присваиванием АЙДИ каждому запросу, что может помочь при отслеживании
	// действий для конкретного запроса (например при ошибке, чтоб понять, где ошибка)
	router.Use(middleware.Logger) //Свой встроенный логгер в мидлвейр от CHI, будет красиво логгировать каждый запрос клиента, со временем на обработку и т.д.
	router.Use(mwLogger.New(logger))
	router.Use(middleware.Recoverer) // в случае паники восстанавливаем, чтоб ничего не падало из-за 1 запроса
	router.Use(middleware.URLFormat)
	//Авторизация Админа
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.UrlSaveHandler(logger, storage, cfg))
		r.Delete("/{alias}", deleteReq.UrlDeleteHandler(logger, storage))
	})
	router.Get("/{alias}", redirect.UrlRedirectHandler(logger, storage))
	return router
}

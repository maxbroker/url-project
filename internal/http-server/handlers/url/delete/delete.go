package delete

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/storage"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLDeleter interface {
	DeleteUrl(alias string) error
}

func UrlDeleteHandler(logger *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.delete.UrlDeleteHandler"

		logger := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		alias := chi.URLParam(request, "alias")
		if alias == "" {
			logger.Info("alias is empty")
			render.JSON(writer, request, resp.Error("invalid request"))
			return
		}

		err := urlDeleter.DeleteUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found", slog.String("alias", alias))
			render.JSON(writer, request, resp.Error(fmt.Sprintf("URL with alias '%s' not found", alias)))
			return
		}
		if err != nil {
			logger.Error("failed to delete url", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to delete url"))
			return
		}
		logger.Info("url deleted", slog.String("alias", alias))
		responseOK(writer, request, alias)
	}
}

func responseOK(writer http.ResponseWriter, request *http.Request, alias string) {
	render.JSON(writer, request, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

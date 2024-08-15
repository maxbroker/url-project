package redirect

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type UrlGetter interface {
	GetUrl(alias string) (string, error)
}

func RedirectUrl(logger *slog.Logger, urlGetter UrlGetter) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.redirect.RedirectUrl"

		logger = logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)
		alias := chi.URLParam(request, "alias")
		if alias == "" {
			logger.Info("alias is empty")

			render.JSON(writer, request, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetUrl(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			logger.Info("url not found", "alias", alias)

			render.JSON(writer, request, resp.Error("not found"))

			return
		}
		if err != nil {
			logger.Error("failed to get url", sl.Err(err))

			render.JSON(writer, request, resp.Error("internal error"))

			return
		}
		logger.Info("got url", slog.String("url", resURL))
		http.Redirect(writer, request, resURL, http.StatusFound)
	}
}

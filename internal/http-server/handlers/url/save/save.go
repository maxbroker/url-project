package save

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-playground/validator/v10"
	"log/slog"

	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/logger/sl"
	"awesomeProject/internal/lib/random"
	"awesomeProject/internal/storage"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config if needed
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.44.1 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (primitive.ObjectID, error)
}

func SaveUrlHandler(logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		const op = "handlers.url.save.SaveUrlHandler"

		logger := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(request.Context())),
		)

		var req Request
		err := render.DecodeJSON(request.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			logger.Error("request body is empty")
			render.JSON(writer, request, resp.Error("empty request"))
			return
		}
		if err != nil {
			logger.Error("failed to decode request body", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to decode request"))
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			logger.Error("invalid request", sl.Err(err))
			render.JSON(writer, request, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			logger.Info("url already exists", slog.String("url", req.URL))
			render.JSON(writer, request, resp.Error("url already exists"))
			return
		}
		if err != nil {
			logger.Error("failed to add url", sl.Err(err))
			render.JSON(writer, request, resp.Error("failed to add url"))
			return
		}

		logger.Info("url added", slog.Any("id", id))
		responseOK(writer, request, alias)
	}
}

func responseOK(writer http.ResponseWriter, request *http.Request, alias string) {
	render.JSON(writer, request, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}

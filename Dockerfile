FROM golang:1.22.6-alpine AS builder

RUN apk --no-cache add bash git make gcc gettext musl-dev

# Обновляем путь для go.mod и go.sum
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Копируем оставшиеся файлы проекта
COPY . ./
RUN go build -o /bin/app cmd/url-shortener/main.go

FROM alpine AS runner

COPY --from=builder /bin/app /
COPY config/config.yaml /config.yaml
COPY db/migrations /db/migrations

CMD ["/app"]

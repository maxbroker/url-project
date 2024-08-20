FROM golang:1.22.6-alpine AS builder

RUN apk --no-cache add bash git make gcc gettext musl-dev

WORKDIR /usr/src


COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download


COPY app ./
RUN go build -o ./bin/app cmd/url-shortener/main.go

FROM alpine AS runner

COPY --from=builder /usr/src/bin/app /
COPY config/config.yaml /config.yaml

CMD ["/app"]
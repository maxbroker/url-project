package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func (d DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (d DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (d DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return d
}

func (d DiscardHandler) WithGroup(name string) slog.Handler {
	return d
}

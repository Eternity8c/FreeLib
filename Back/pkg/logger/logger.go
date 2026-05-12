package logger

import (
	"io"
	"log/slog"
)

func NewLogger(Type string) *slog.Logger {
	switch Type {
	case "JSON":
		return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
	case "Text":
		return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	default:
		return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
	}
}

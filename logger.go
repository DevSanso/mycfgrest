package main

import (
	"io"
	"golang.org/x/exp/slog"
)

func NewLogger(writer io.Writer, level string) *slog.Logger {
	levelVar := new(slog.LevelVar)

	switch level {
	case "debug":
		levelVar.Set(slog.LevelDebug)
	case "error":
		levelVar.Set(slog.LevelError)
	case "warn":
		levelVar.Set(slog.LevelWarn)
	default:
		levelVar.Set(slog.LevelInfo)
	}

	handle := slog.NewTextHandler(writer, &slog.HandlerOptions {
		Level: levelVar,
	})

	return slog.New(handle)	
}
package utils

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {
	Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
}

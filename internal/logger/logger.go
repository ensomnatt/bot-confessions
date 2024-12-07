package logger

import (
  "log/slog"
  "os"
)

var Logger *slog.Logger

func Init(env string) {
  if env == "prod" {
    Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
      Level: slog.LevelInfo,
    }))
  } else if env == "dev" {
    Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
      Level: slog.LevelDebug,
    }))
  } else if env == "local" {
    Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
      Level: slog.LevelDebug,
    }))
  } else {
    Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
      Level: slog.LevelInfo,
    }))
  }
}

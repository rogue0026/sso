package main

import (
	"log/slog"
	"os"

	"github.com/rogue0026/sso/internal/app"
	"github.com/rogue0026/sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// init app config
	appCfg := config.MustLoad()

	// init app logger
	log := setupLogger(appCfg.Env)
	log.Debug("logger initialized")

	application, err := app.New(appCfg, log)
	if err != nil {
		panic(err.Error())
	}

	log.Debug("starting listening at ", "port", appCfg.GRPC.Port)
	application.MustRun(appCfg.GRPC.Port)
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))
	}

	return logger
}

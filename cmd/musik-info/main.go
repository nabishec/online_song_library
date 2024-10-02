package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/nabishec/restapi/internal/config"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/storage/postgresql"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	log.Println("config of system:", cfg)
	// TODO: init logger: slog
	log := setUpLogger(cfg.Env)

	log.Info("hello world", slog.String("env", cfg.Env))
	log.Debug("message")
	// TODO: init storage: postgresql
	storage, err := postgresql.NewDatabase(cfg.DBDataSourceName)
	if err != nil {
		log.Error("failed to init storage", slerr.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO: init router: chi, "chi render"

	//TODO: run server:
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setUpLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

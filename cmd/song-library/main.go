package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nabishec/restapi/internal/config"
	"github.com/nabishec/restapi/internal/http-server/handlers/deletion"
	"github.com/nabishec/restapi/internal/http-server/handlers/get"
	"github.com/nabishec/restapi/internal/http-server/handlers/post"
	"github.com/nabishec/restapi/internal/http-server/handlers/put"
	"github.com/nabishec/restapi/internal/http-server/middleware/logger"
	"github.com/nabishec/restapi/internal/lib/logger/slerr"
	"github.com/nabishec/restapi/internal/storage/postgresql"

	_ "github.com/nabishec/restapi/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Song Library
// @version 1.0
// @description API Server for SongLibrary
// @contact.email nabishec@mail.ru

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("Programm started")

	// TODO: init storage: postgresql
	storage, err := postgresql.NewDatabase()
	if err != nil {
		log.Error("failed to init storage", slerr.Err(err))
		os.Exit(1)
	}
	// _ = storage

	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/v1/songslibrary/song", post.SongPost(log, storage))
	router.Get("/api/v1/songslibrary", get.SongsLibrary(log, storage))
	router.Delete("/api/v1/songslibrary/song", deletion.SongDelete(log, storage))
	router.Get("/api/v1/songslibrary/song", get.TextSongGet(log, storage))
	router.Put("/api/v1/songslibrary/song", put.SongDetail(log, storage))
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Info("starting server", slog.String("address", cfg.Address))

	//TODO: run server:
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", slerr.Err(err))
	}
	log.Error("server stoped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
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

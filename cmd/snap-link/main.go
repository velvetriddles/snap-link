package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/velvetriddles/snap-link/internal/config"
	"github.com/velvetriddles/snap-link/internal/http-server/handlers/redirect"
	"github.com/velvetriddles/snap-link/internal/http-server/handlers/urls/save"
	mwLogger "github.com/velvetriddles/snap-link/internal/http-server/middleware/logger"
	"github.com/velvetriddles/snap-link/internal/lib/logger/handlers/slogpretty"
	"github.com/velvetriddles/snap-link/internal/lib/logger/sl"
	"github.com/velvetriddles/snap-link/internal/storage/sqlite"
	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := setupLogger(cfg.Env)
	log.Info("Starting snap-link service", slog.String("env", cfg.Env))
	log.Debug("Debugging snap-link service")
	log.Error("error messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)

	// fmt.Println(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("snap-link", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage)) // Переместили маршрут внутрь подмаршрута
	})
	router.Get("/{alias}", redirect.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

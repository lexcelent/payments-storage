package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/lexcelent/payments-storage/internal/config"
	"github.com/lexcelent/payments-storage/internal/storage/sqlite"
	"github.com/lexcelent/payments-storage/internal/transport/http/handlers"
	"github.com/lexcelent/payments-storage/internal/transport/http/router"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)
	log.Info("setup logger")

	router := router.New()
	log.Debug("setup router")

	router.Handle("/healthCheck", handlers.HealthCheck)
	router.Handle("/payments/add", handlers.PaymentAdd)

	log.Info(
		"server has been started",
		slog.String("address", cfg.HTTPServer.Address),
		slog.String("port", cfg.HTTPServer.Port),
	)

	// TODO: move storage logic
	_, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(err)
	}
	// END TODO

	if err := http.ListenAndServe(":"+cfg.HTTPServer.Port, router); err != nil {
		log.Error("Ошибка запуска HTTP-сервера")
	}
}

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

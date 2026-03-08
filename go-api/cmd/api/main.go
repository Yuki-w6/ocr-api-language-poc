package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/config"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/db"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/repository"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/router"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/service"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	ocrJobRepository := repository.NewOCRJobRepository(pool)
	uploadService := service.NewUploadService(cfg.UploadURLBase, cfg.PresignedURLExpiresIn)
	ocrJobService := service.NewOCRJobService(ocrJobRepository, logger)

	r := router.NewRouter(logger, uploadService, ocrJobService)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("server starting", slog.String("port", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server stopped")
}

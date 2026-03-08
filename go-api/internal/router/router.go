package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/handler"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	logger *slog.Logger,
	uploadService *service.UploadService,
	ocrJobService *service.OCRJobService,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(requestLogger(logger))

	healthHandler := handler.NewHealthHandler()
	ocrJobHandler := handler.NewOCRJobHandler(logger, uploadService, ocrJobService)

	r.Get("/health", healthHandler.GetHealth)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/uploads/presigned-url", ocrJobHandler.CreatePresignedURL)
		r.Post("/ocr-jobs", ocrJobHandler.CreateOCRJob)
		r.Get("/ocr-jobs/{jobId}", ocrJobHandler.GetOCRJobStatus)
		r.Get("/ocr-jobs/{jobId}/result", ocrJobHandler.GetOCRJobResult)
	})

	return r
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger.Info("http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

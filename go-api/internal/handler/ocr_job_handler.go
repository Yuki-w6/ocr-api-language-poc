package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/request"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/response"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type OCRJobHandler struct {
	logger        *slog.Logger
	validate      *validator.Validate
	uploadService *service.UploadService
	ocrJobService *service.OCRJobService
}

func NewOCRJobHandler(
	logger *slog.Logger,
	uploadService *service.UploadService,
	ocrJobService *service.OCRJobService,
) *OCRJobHandler {
	return &OCRJobHandler{
		logger:        logger,
		validate:      validator.New(),
		uploadService: uploadService,
		ocrJobService: ocrJobService,
	}
}

func (h *OCRJobHandler) CreatePresignedURL(w http.ResponseWriter, r *http.Request) {
	var req request.CreatePresignedURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "invalid json body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	result := h.uploadService.CreatePresignedURL(req.Filename)
	response.JSON(w, http.StatusOK, result)
}

func (h *OCRJobHandler) CreateOCRJob(w http.ResponseWriter, r *http.Request) {
	var req request.CreateOCRJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", "invalid json body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	result, err := h.ocrJobService.CreateJob(r.Context(), req.ObjectKey)
	if err != nil {
		h.logger.Error("failed to create ocr job", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "internal_error", "failed to create ocr job")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *OCRJobHandler) GetOCRJobStatus(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobId")
	if jobID == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "jobId is required")
		return
	}

	result, err := h.ocrJobService.GetJobStatus(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, service.ErrJobNotFound) {
			response.Error(w, http.StatusNotFound, "not_found", "ocr job not found")
			return
		}

		h.logger.Error("failed to get ocr job status", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "internal_error", "failed to get ocr job status")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *OCRJobHandler) GetOCRJobResult(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "jobId")
	if jobID == "" {
		response.Error(w, http.StatusBadRequest, "invalid_request", "jobId is required")
		return
	}

	result, err := h.ocrJobService.GetJobResult(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, service.ErrJobNotFound) {
			response.Error(w, http.StatusNotFound, "not_found", "ocr job not found")
			return
		}

		h.logger.Error("failed to get ocr job result", slog.String("error", err.Error()))
		response.Error(w, http.StatusInternalServerError, "internal_error", "failed to get ocr job result")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

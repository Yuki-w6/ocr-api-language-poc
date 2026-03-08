package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/model"
	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/service"
	"github.com/go-chi/chi/v5"
)

type fakeOCRJobRepository struct {
	createFn   func(ctx context.Context, job *model.OCRJob) error
	findByIDFn func(ctx context.Context, id string) (*model.OCRJob, error)
}

func (f *fakeOCRJobRepository) Create(ctx context.Context, job *model.OCRJob) error {
	if f.createFn != nil {
		return f.createFn(ctx, job)
	}
	return nil
}

func (f *fakeOCRJobRepository) FindByID(ctx context.Context, id string) (*model.OCRJob, error) {
	if f.findByIDFn != nil {
		return f.findByIDFn(ctx, id)
	}
	return nil, nil
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

func newTestHandler(repo service.OCRJobRepository) *OCRJobHandler {
	logger := newTestLogger()
	uploadService := service.NewUploadService("https://storage.example.com/upload", 300)
	ocrJobService := service.NewOCRJobService(repo, logger)

	return NewOCRJobHandler(logger, uploadService, ocrJobService)
}

func TestOCRJobHandler_CreatePresignedURL(t *testing.T) {
	t.Parallel()

	h := newTestHandler(&fakeOCRJobRepository{})

	body := `{
		"filename": "receipt.jpg",
		"contentType": "image/jpeg"
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/uploads/presigned-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreatePresignedURL(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp struct {
		Data struct {
			ObjectKey string `json:"objectKey"`
			UploadURL string `json:"uploadUrl"`
			ExpiresIn int    `json:"expiresIn"`
		} `json:"data"`
	}
	if err := json.NewDecoder(bytes.NewReader(rec.Body.Bytes())).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Data.ObjectKey == "" {
		t.Fatal("expected objectKey to be set")
	}
	if resp.Data.UploadURL == "" {
		t.Fatal("expected uploadUrl to be set")
	}
	if resp.Data.ExpiresIn != 300 {
		t.Fatalf("expected expiresIn=300, got %d", resp.Data.ExpiresIn)
	}
}

func TestOCRJobHandler_CreatePresignedURL_ValidationError(t *testing.T) {
	t.Parallel()

	h := newTestHandler(&fakeOCRJobRepository{})

	body := `{
		"filename": "",
		"contentType": "image/jpeg"
	}`

	req := httptest.NewRequest(http.MethodPost, "/v1/uploads/presigned-url", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreatePresignedURL(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestOCRJobHandler_GetOCRJobStatus_NotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeOCRJobRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.OCRJob, error) {
			return nil, nil
		},
	}
	h := newTestHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/v1/ocr-jobs/job_not_found", nil)
	rec := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("jobId", "job_not_found")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	h.GetOCRJobStatus(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	var resp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(bytes.NewReader(rec.Body.Bytes())).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Error.Code != "not_found" {
		t.Fatalf("expected error code not_found, got %s", resp.Error.Code)
	}
}

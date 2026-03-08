package service

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/model"
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

func TestOCRJobService_CreateJob(t *testing.T) {
	t.Parallel()

	var createdJob *model.OCRJob

	repo := &fakeOCRJobRepository{
		createFn: func(ctx context.Context, job *model.OCRJob) error {
			createdJob = job
			return nil
		},
	}

	svc := NewOCRJobService(repo, newTestLogger())

	result, err := svc.CreateJob(context.Background(), "uploads/2026/03/08/test.jpg")
	if err != nil {
		t.Fatalf("CreateJob() error = %v", err)
	}

	if result.JobID == "" {
		t.Fatal("expected job id to be set")
	}
	if result.Status != model.JobStatusQueued {
		t.Fatalf("expected status %q, got %q", model.JobStatusQueued, result.Status)
	}

	if createdJob == nil {
		t.Fatal("expected repository Create to be called")
	}
	if createdJob.ObjectKey != "uploads/2026/03/08/test.jpg" {
		t.Fatalf("unexpected object key: %s", createdJob.ObjectKey)
	}
	if createdJob.Status != model.JobStatusQueued {
		t.Fatalf("unexpected created job status: %s", createdJob.Status)
	}
}

func TestOCRJobService_GetJobStatus_NotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeOCRJobRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.OCRJob, error) {
			return nil, nil
		},
	}

	svc := NewOCRJobService(repo, newTestLogger())

	_, err := svc.GetJobStatus(context.Background(), "job_not_found")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err != ErrJobNotFound {
		t.Fatalf("expected ErrJobNotFound, got %v", err)
	}
}

func TestOCRJobService_GetJobResult(t *testing.T) {
	t.Parallel()

	repo := &fakeOCRJobRepository{
		findByIDFn: func(ctx context.Context, id string) (*model.OCRJob, error) {
			return &model.OCRJob{
				ID:         "job_123",
				ObjectKey:  "uploads/2026/03/08/test.jpg",
				Status:     model.JobStatusSucceeded,
				ResultJSON: []byte(`{"text":"sample OCR result"}`),
			}, nil
		},
	}

	svc := NewOCRJobService(repo, newTestLogger())

	result, err := svc.GetJobResult(context.Background(), "job_123")
	if err != nil {
		t.Fatalf("GetJobResult() error = %v", err)
	}

	if result.JobID != "job_123" {
		t.Fatalf("unexpected job id: %s", result.JobID)
	}
	if result.Status != model.JobStatusSucceeded {
		t.Fatalf("unexpected status: %s", result.Status)
	}
	if string(result.Result) != `{"text":"sample OCR result"}` {
		t.Fatalf("unexpected result json: %s", string(result.Result))
	}
}

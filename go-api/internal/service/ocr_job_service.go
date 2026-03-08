package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/model"
	"github.com/google/uuid"
)

var ErrJobNotFound = errors.New("job not found")

type OCRJobRepository interface {
	Create(ctx context.Context, job *model.OCRJob) error
	FindByID(ctx context.Context, id string) (*model.OCRJob, error)
}

type OCRJobService struct {
	repo   OCRJobRepository
	logger *slog.Logger
}

type CreateOCRJobResult struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
}

type OCRJobStatusResult struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
}

type OCRJobResultResponse struct {
	JobID  string          `json:"jobId"`
	Status string          `json:"status"`
	Result json.RawMessage `json:"result,omitempty"`
}

func NewOCRJobService(repo OCRJobRepository, logger *slog.Logger) *OCRJobService {
	return &OCRJobService{
		repo:   repo,
		logger: logger,
	}
}

func (s *OCRJobService) CreateJob(ctx context.Context, objectKey string) (*CreateOCRJobResult, error) {
	jobID := "job_" + uuid.NewString()

	job := &model.OCRJob{
		ID:         jobID,
		ObjectKey:  objectKey,
		Status:     model.JobStatusQueued,
		ResultJSON: nil,
	}

	if err := s.repo.Create(ctx, job); err != nil {
		return nil, err
	}

	s.logger.Info("ocr job queued",
		slog.String("job_id", jobID),
		slog.String("object_key", objectKey),
	)

	// PoC では Queue 送信の代わりにログ出力のみ
	s.logger.Info("queue publish simulated",
		slog.String("job_id", jobID),
		slog.String("object_key", objectKey),
	)

	return &CreateOCRJobResult{
		JobID:  jobID,
		Status: job.Status,
	}, nil
}

func (s *OCRJobService) GetJobStatus(ctx context.Context, jobID string) (*OCRJobStatusResult, error) {
	job, err := s.repo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	return &OCRJobStatusResult{
		JobID:  job.ID,
		Status: job.Status,
	}, nil
}

func (s *OCRJobService) GetJobResult(ctx context.Context, jobID string) (*OCRJobResultResponse, error) {
	job, err := s.repo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrJobNotFound
	}

	resp := &OCRJobResultResponse{
		JobID:  job.ID,
		Status: job.Status,
	}

	if len(job.ResultJSON) > 0 {
		resp.Result = json.RawMessage(job.ResultJSON)
	}

	return resp, nil
}

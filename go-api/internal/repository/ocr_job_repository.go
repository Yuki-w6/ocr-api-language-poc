package repository

import (
	"context"
	"errors"

	"github.com/Yuki-w6/ocr-api-language-poc/go-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OCRJobRepository struct {
	pool *pgxpool.Pool
}

func NewOCRJobRepository(pool *pgxpool.Pool) *OCRJobRepository {
	return &OCRJobRepository{pool: pool}
}

func (r *OCRJobRepository) Create(ctx context.Context, job *model.OCRJob) error {
	const query = `
		INSERT INTO ocr_jobs (id, object_key, status, result_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := r.pool.Exec(ctx, query, job.ID, job.ObjectKey, job.Status, job.ResultJSON)
	return err
}

func (r *OCRJobRepository) FindByID(ctx context.Context, id string) (*model.OCRJob, error) {
	const query = `
		SELECT id, object_key, status, result_json, created_at, updated_at
		FROM ocr_jobs
		WHERE id = $1
	`

	var job model.OCRJob
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.ObjectKey,
		&job.Status,
		&job.ResultJSON,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &job, nil
}

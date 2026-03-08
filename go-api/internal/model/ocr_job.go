package model

import "time"

const (
	JobStatusQueued     = "queued"
	JobStatusProcessing = "processing"
	JobStatusSucceeded  = "succeeded"
	JobStatusFailed     = "failed"
)

type OCRJob struct {
	ID         string
	ObjectKey  string
	Status     string
	ResultJSON []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

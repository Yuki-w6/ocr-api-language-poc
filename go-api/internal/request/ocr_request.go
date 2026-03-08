package request

type CreatePresignedURLRequest struct {
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"contentType" validate:"required"`
}

type CreateOCRJobRequest struct {
	ObjectKey string `json:"objectKey" validate:"required"`
}

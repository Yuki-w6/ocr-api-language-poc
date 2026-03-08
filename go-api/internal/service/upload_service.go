package service

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UploadService struct {
	uploadURLBase string
	expiresIn     int
}

type PresignedUploadResult struct {
	ObjectKey string `json:"objectKey"`
	UploadURL string `json:"uploadUrl"`
	ExpiresIn int    `json:"expiresIn"`
}

func NewUploadService(uploadURLBase string, expiresIn int) *UploadService {
	return &UploadService{
		uploadURLBase: uploadURLBase,
		expiresIn:     expiresIn,
	}
}

func (s *UploadService) CreatePresignedURL(filename string) PresignedUploadResult {
	now := time.Now().UTC()
	safeFilename := sanitizeFilename(filename)
	objectKey := fmt.Sprintf(
		"uploads/%04d/%02d/%02d/%s-%s",
		now.Year(),
		now.Month(),
		now.Day(),
		uuid.NewString(),
		safeFilename,
	)

	return PresignedUploadResult{
		ObjectKey: objectKey,
		UploadURL: fmt.Sprintf("%s/%s", strings.TrimRight(s.uploadURLBase, "/"), objectKey),
		ExpiresIn: s.expiresIn,
	}
}

func sanitizeFilename(filename string) string {
	base := filepath.Base(strings.TrimSpace(filename))
	if base == "." || base == "" || base == "/" {
		return "file.bin"
	}
	return base
}

package storage

import (
	"context"
	"io"
)

// Uploader adalah interface standar untuk semua layanan penyimpanan objek.
type Uploader interface {
	UploadFile(ctx context.Context, file io.Reader, objectName string, contentType string, size int64) (string, error)
	DeleteFile(ctx context.Context, objectName string) error
	GetBaseURL() string
}

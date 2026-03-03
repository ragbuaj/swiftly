package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioUploader struct {
	client     *minio.Client
	bucketName string
	endpoint   string // Internal endpoint (e.g. minio:9000)
	publicURL  string // External endpoint (e.g. http://localhost:9000)
	useSSL     bool
}

// NewMinioUploader menginisialisasi koneksi ke server MinIO / S3.
func NewMinioUploader(endpoint, accessKey, secretKey, bucketName, publicURL string, useSSL bool) (*MinioUploader, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %v", err)
	}

	log.Printf("Successfully connected to MinIO/S3 at %s (Public: %s)", endpoint, publicURL)

	return &MinioUploader{
		client:     minioClient,
		bucketName: bucketName,
		endpoint:   endpoint,
		publicURL:  publicURL,
		useSSL:     useSSL,
	}, nil
}

func (m *MinioUploader) GetBaseURL() string {
	if m.publicURL != "" {
		return fmt.Sprintf("%s/%s", m.publicURL, m.bucketName)
	}

	protocol := "http"
	if m.useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s", protocol, m.endpoint, m.bucketName)
}

// UploadFile mengunggah file ke MinIO dan mengembalikan object name (key) saja.
func (m *MinioUploader) UploadFile(ctx context.Context, file io.Reader, objectName string, contentType string, size int64) (string, error) {
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectName, file, size, opts)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to minio: %v", err)
	}

	// Hanya kembalikan objectName (Key), bukan URL lengkap.
	return objectName, nil
}

// DeleteFile menghapus file dari bucket.
func (m *MinioUploader) DeleteFile(ctx context.Context, objectName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, opts)
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %v", err)
	}
	return nil
}

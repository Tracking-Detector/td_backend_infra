package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type MinIOStorageAdapter struct {
	client *minio.Client
}

func NewMinIOStorageAdapter(client *minio.Client) *MinIOStorageAdapter {
	return &MinIOStorageAdapter{
		client: client,
	}

}

func (a *MinIOStorageAdapter) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return a.client.BucketExists(ctx, bucketName)
}
func (a *MinIOStorageAdapter) MakeBucket(ctx context.Context, bucketName string, options minio.MakeBucketOptions) error {
	return a.MakeBucket(ctx, bucketName, options)
}
func (a *MinIOStorageAdapter) SetBucketVersioning(ctx context.Context, bucketName string, config minio.BucketVersioningConfiguration) error {
	return a.SetBucketVersioning(ctx, bucketName, config)
}
func (a *MinIOStorageAdapter) GetObject(ctx context.Context, bucketName string, fileURI string, options minio.GetObjectOptions) (io.ReadSeekCloser, error) {
	return a.GetObject(ctx, bucketName, fileURI, options)
}
func (a *MinIOStorageAdapter) PutObject(ctx context.Context, bucketName, fileName string, pr io.Reader, objectSize int64, options minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	return a.PutObject(ctx, bucketName, fileName, pr, objectSize, options)
}
func (a *MinIOStorageAdapter) ListObjects(ctx context.Context, bucketName string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	return a.ListObjects(ctx, bucketName, options)
}

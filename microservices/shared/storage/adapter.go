package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type IStorageClientAdater interface {
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	MakeBucket(ctx context.Context, bucketName string, options minio.MakeBucketOptions) error
	SetBucketVersioning(ctx context.Context, bucketName string, config minio.BucketVersioningConfiguration) error
	GetObject(ctx context.Context, bucketName string, fileURI string, options minio.GetObjectOptions) (io.ReadSeekCloser, error)
	PutObject(ctx context.Context, bucketName, fileName string, pr io.Reader, objectSize int64, options minio.PutObjectOptions) (info minio.UploadInfo, err error)
	ListObjects(ctx context.Context, bucketName string, options minio.ListObjectsOptions) <-chan minio.ObjectInfo
}

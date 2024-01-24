package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"tds/shared/storage"

	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

type IStorageService interface {
	VerifyBucketExists(ctx context.Context, bucketName string)
	DownloadFile(ctx context.Context, bucketName, fileURI, destination string) error
	GetObject(ctx context.Context, bucketName string, filename string) (io.ReadSeekCloser, error)
	PutObject(ctx context.Context, bucketName string, fileName string, pr *io.PipeReader, objectSize int64, contentType string) error
	GetBucketStructure(bucketName, prefix string) (interface{}, error)
}

type MinIOStorageService struct {
	client storage.IStorageClientAdater
}

func NewMinIOStorageService(client storage.IStorageClientAdater) *MinIOStorageService {
	return &MinIOStorageService{
		client: client,
	}
}

func (s *MinIOStorageService) VerifyBucketExists(ctx context.Context, bucketName string) {
	if exists, err := s.client.BucketExists(ctx, bucketName); err != nil {
		log.WithFields(log.Fields{
			"service": "setup",
			"error":   err.Error(),
		}).Fatal("Error verifing whether bucket exisits.")
	} else if exists {
	} else {
		if makeBucketError := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "eu-central-1"}); makeBucketError != nil {
			log.WithFields(log.Fields{
				"service": "setup",
				"error":   makeBucketError.Error(),
			}).Fatal("Error creating bucket with name ", bucketName, ".")
		} else {
			if setVersioningError := s.client.SetBucketVersioning(ctx, bucketName, minio.BucketVersioningConfiguration{
				Status: "Enabled",
			}); setVersioningError != nil {
				log.WithFields(log.Fields{
					"service": "setup",
					"error":   makeBucketError.Error(),
				}).Fatal("Error setting versioning for bucket with name ", bucketName, ".")
			}
		}
	}
}

func (s *MinIOStorageService) DownloadFile(ctx context.Context, bucketName, fileURI, destination string) error {
	obj, err := s.client.GetObject(ctx, bucketName, fileURI, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer obj.Close()
	fileExt := filepath.Ext(fileURI)
	file, err := os.Create(destination + fileExt)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, obj)
	if err != nil {
		return err
	}

	return nil
}

func (s *MinIOStorageService) PutObject(ctx context.Context, bucketName string, fileName string, pr *io.PipeReader, objectSize int64, contentType string) error {
	_, err := s.client.PutObject(ctx, bucketName, fileName, pr, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *MinIOStorageService) GetObject(ctx context.Context, bucketName string, filename string) (io.ReadSeekCloser, error) {
	return s.client.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
}

func (s *MinIOStorageService) GetBucketStructure(bucketName, prefix string) (interface{}, error) {
	bucketStructure := make(map[string]interface{})

	doneCh := make(chan struct{})
	defer close(doneCh)

	// Retrieve all objects in the specified bucket and prefix
	objectsCh := s.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Recursive:    true,
		WithVersions: true,
		WithMetadata: true,
	})
	for object := range objectsCh {
		if object.Err != nil {
			return nil, object.Err
		}

		// Split the object key into directories and filename
		directories := s.splitDirectories(object.Key)
		// Build the nested structure based on the directories
		currentDir := bucketStructure
		for _, directory := range directories {
			if _, ok := currentDir[directory]; !ok {
				currentDir[directory] = make(map[string]interface{})
			}
			currentDir = currentDir[directory].(map[string]interface{})
		}
	}

	return bucketStructure, nil
}

func (s *MinIOStorageService) splitDirectories(key string) []string {
	directories := make([]string, 0)

	// Remove leading and trailing slashes
	key = s.trimSlashes(key)

	// Split the key into directories
	for _, dir := range strings.Split(key, "/") {
		if dir != "" {
			directories = append(directories, dir)
		}
	}

	return directories
}

// Helper function to remove leading and trailing slashes from a string
func (s *MinIOStorageService) trimSlashes(str string) string {
	return strings.Trim(str, "/")
}

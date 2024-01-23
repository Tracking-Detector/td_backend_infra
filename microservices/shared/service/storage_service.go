package service

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"tds/shared/configs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type IStorageService interface {
}

type MinIOStorageService struct {
	client *minio.Client
}

func NewMinIOStorageServic() *MinIOStorageService {
	minioClient, err := minio.New(configs.EnvMinIoURI(), &minio.Options{
		Creds:  credentials.NewStaticV4(configs.EnvMinIoAccessKey(), configs.EnvMinIoPrivateKey(), ""),
		Secure: false,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"service": "MinIOStorageService",
		}).Fatalln(err)
	}
	log.WithFields(log.Fields{
		"service": "MinIOStorageService",
	}).Info("Successfully connected to MinIO.")
	return &MinIOStorageService{
		client: minioClient,
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

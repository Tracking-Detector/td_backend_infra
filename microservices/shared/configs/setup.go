package configs

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{
		"service": "setup",
	}).Info("Successfully connected to MongoDB.")
	return client
}

func ConnectMinio() *minio.Client {
	minioClient, err := minio.New(EnvMinIoURI(), &minio.Options{
		Creds:  credentials.NewStaticV4(EnvMinIoAccessKey(), EnvMinIoPrivateKey(), ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.WithFields(log.Fields{
		"service": "setup",
	}).Info("Successfully connected to MinIO.")
	return minioClient
}

var DB *mongo.Client = ConnectDB()

var MINIO *minio.Client = ConnectMinio()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("tracking-detector").Collection(collectionName)
	return collection
}

func VerifyBucketExists(ctx context.Context, client *minio.Client, bucketName string) {
	if exists, err := client.BucketExists(ctx, bucketName); err != nil {
		log.WithFields(log.Fields{
			"service": "setup",
			"error":   err.Error(),
		}).Fatal("Error verifing whether bucket exisits.")
	} else if exists {
	} else {
		if makeBucketError := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "eu-central-1"}); makeBucketError != nil {
			log.WithFields(log.Fields{
				"service": "setup",
				"error":   makeBucketError.Error(),
			}).Fatal("Error creating bucket with name ", bucketName, ".")
		} else {
			if setVersioningError := client.SetBucketVersioning(ctx, bucketName, minio.BucketVersioningConfiguration{
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

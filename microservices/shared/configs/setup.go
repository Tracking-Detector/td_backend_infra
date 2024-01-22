package configs

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context) *mongo.Client {
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

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("tracking-detector").Collection(collectionName)
	return collection
}

func GetDatabase(client *mongo.Client) *mongo.Database {
	db := client.Database("tracking-detector")
	return db
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

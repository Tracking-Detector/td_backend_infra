package utils

import (
	"context"
	"strings"

	"github.com/minio/minio-go/v7"
)

func GetBucketStructure(client *minio.Client, bucketName, prefix string) (interface{}, error) {
	bucketStructure := make(map[string]interface{})

	doneCh := make(chan struct{})
	defer close(doneCh)

	// Retrieve all objects in the specified bucket and prefix
	objectsCh := client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Recursive:    true,
		WithVersions: true,
		WithMetadata: true,
	})
	for object := range objectsCh {
		if object.Err != nil {
			return nil, object.Err
		}

		// Split the object key into directories and filename
		directories := SplitDirectories(object.Key)
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

func SplitDirectories(key string) []string {
	directories := make([]string, 0)

	// Remove leading and trailing slashes
	key = TrimSlashes(key)

	// Split the key into directories
	for _, dir := range strings.Split(key, "/") {
		if dir != "" {
			directories = append(directories, dir)
		}
	}

	return directories
}

// Helper function to remove leading and trailing slashes from a string
func TrimSlashes(str string) string {
	return strings.Trim(str, "/")
}

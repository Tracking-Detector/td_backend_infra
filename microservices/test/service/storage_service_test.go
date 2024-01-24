package service_test

import (
	"context"
	"errors"
	"io"
	"os"
	"tds/mocks"
	"tds/shared/service"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestStorageService(t *testing.T) {
	suite.Run(t, &StorageServiceTest{})
}

type StorageServiceTest struct {
	suite.Suite
	storageService *service.MinIOStorageService
	minioAdapter   *mocks.IStorageClientAdater
}

func (suite *StorageServiceTest) SetupTest() {
	suite.minioAdapter = mocks.NewIStorageClientAdater(suite.T())
	suite.storageService = service.NewMinIOStorageService(suite.minioAdapter)
}

func (suite *StorageServiceTest) TestVerifyBucketExists_BucketAlreadyExists() {
	// given
	suite.minioAdapter.On("BucketExists", mock.Anything, "existingBucket").Return(true, nil)

	// when
	suite.storageService.VerifyBucketExists(context.Background(), "existingBucket")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "BucketExists", mock.Anything, "existingBucket")
}

func (suite *StorageServiceTest) TestVerifyBucketExists_CreateBucketSuccess() {
	// given
	suite.minioAdapter.On("BucketExists", mock.Anything, "newBucket").Return(false, nil)
	suite.minioAdapter.On("MakeBucket", mock.Anything, "newBucket", mock.AnythingOfType("minio.MakeBucketOptions")).Return(nil)
	suite.minioAdapter.On("SetBucketVersioning", mock.Anything, "newBucket", mock.AnythingOfType("minio.BucketVersioningConfiguration")).Return(nil)

	// when
	suite.storageService.VerifyBucketExists(context.Background(), "newBucket")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "BucketExists", mock.Anything, "newBucket")
	suite.minioAdapter.AssertCalled(suite.T(), "MakeBucket", mock.Anything, "newBucket", mock.AnythingOfType("minio.MakeBucketOptions"))
	suite.minioAdapter.AssertCalled(suite.T(), "SetBucketVersioning", mock.Anything, "newBucket", mock.AnythingOfType("minio.BucketVersioningConfiguration"))
}

func (suite *StorageServiceTest) TestDownloadFile_Success() {
	// given
	fileUri := "exporter.js"
	file, err := os.Open("../resources/example.js")
	suite.minioAdapter.On("GetObject", mock.Anything, "bucketName", fileUri, mock.AnythingOfType("minio.GetObjectOptions")).Return(file, nil)

	// when
	err = suite.storageService.DownloadFile(context.Background(), "bucketName", fileUri, "destination")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "GetObject", mock.Anything, "bucketName", fileUri, mock.AnythingOfType("minio.GetObjectOptions"))
	suite.NoError(err)
	os.Remove("./destination.js")
}

func (suite *StorageServiceTest) TestDownloadFile_Error() {
	// given
	suite.minioAdapter.On("GetObject", mock.Anything, "bucketName", "fileURI", mock.AnythingOfType("minio.GetObjectOptions")).Return(nil, errors.New("Download error"))

	// when
	err := suite.storageService.DownloadFile(context.Background(), "bucketName", "fileURI", "destination")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "GetObject", mock.Anything, "bucketName", "fileURI", mock.AnythingOfType("minio.GetObjectOptions"))
	suite.Error(err)
}

func (suite *StorageServiceTest) TestPutObject_Success() {
	// given
	suite.minioAdapter.On("PutObject", mock.Anything, "bucketName", "fileName", mock.Anything, int64(0), minio.PutObjectOptions{
		ContentType: "contentType",
	}).Return(minio.UploadInfo{}, nil)

	// when
	err := suite.storageService.PutObject(context.Background(), "bucketName", "fileName", &io.PipeReader{}, int64(0), "contentType")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "PutObject", mock.Anything, "bucketName", "fileName", mock.Anything, int64(0), minio.PutObjectOptions{
		ContentType: "contentType",
	})
	suite.NoError(err)
}

func (suite *StorageServiceTest) TestPutObject_Error() {
	// given
	suite.minioAdapter.On("PutObject", mock.Anything, "bucketName", "fileName", mock.Anything, int64(0), mock.AnythingOfType("minio.PutObjectOptions")).Return(minio.UploadInfo{}, errors.New("Error"))

	// when
	err := suite.storageService.PutObject(context.Background(), "bucketName", "fileName", &io.PipeReader{}, int64(0), "contentType")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "PutObject", mock.Anything, "bucketName", "fileName", mock.Anything, int64(0), mock.AnythingOfType("minio.PutObjectOptions"))
	suite.Error(err)
}

func (suite *StorageServiceTest) TestGetObject_Success() {
	// given
	fileName := "example.js"
	file, err := os.Open("../resources/example.js")
	suite.minioAdapter.On("GetObject", mock.Anything, "bucketName", fileName, mock.AnythingOfType("minio.GetObjectOptions")).Return(file, nil)

	// when
	result, err := suite.storageService.GetObject(context.Background(), "bucketName", fileName)

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "GetObject", mock.Anything, "bucketName", fileName, mock.AnythingOfType("minio.GetObjectOptions"))
	suite.NoError(err)

	// Validate that result is a non-nil implementation of io.ReadSeekCloser
	suite.NotNil(result)
	_, isReadSeekCloser := result.(io.ReadSeekCloser)
	suite.True(isReadSeekCloser)

	// Cleanup: Close the file
	file.Close()
}

func (suite *StorageServiceTest) TestGetObject_Error() {
	// given
	suite.minioAdapter.On("GetObject", mock.Anything, "bucketName", "fileURI", mock.AnythingOfType("minio.GetObjectOptions")).Return(nil, errors.New("Get object error"))

	// when
	result, err := suite.storageService.GetObject(context.Background(), "bucketName", "fileURI")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "GetObject", mock.Anything, "bucketName", "fileURI", mock.AnythingOfType("minio.GetObjectOptions"))
	suite.Error(err)
	suite.Nil(result)
}

func (suite *StorageServiceTest) TestGetBucketStructure_Success() {
	// given
	objects := []minio.ObjectInfo{
		{Key: "dir1/file1.txt"},
		{Key: "dir2/subdir/file2.txt"},
		{Key: "dir1/file3.txt"},
	}

	suite.minioAdapter.On("ListObjects", mock.Anything, "bucketName", mock.AnythingOfType("minio.ListObjectsOptions")).Return(suite.createChan(objects))

	// when
	result, err := suite.storageService.GetBucketStructure("bucketName", "")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "ListObjects", mock.Anything, "bucketName", mock.AnythingOfType("minio.ListObjectsOptions"))
	suite.NoError(err)

	// Validate the structure
	suite.Equal(map[string]interface{}{
		"dir1": map[string]interface{}{
			"file1.txt": map[string]interface{}{},
			"file3.txt": map[string]interface{}{},
		},
		"dir2": map[string]interface{}{
			"subdir": map[string]interface{}{
				"file2.txt": map[string]interface{}{},
			},
		},
	}, result)
}

func (suite *StorageServiceTest) TestGetBucketStructure_ListObjectsError() {
	// given
	objects := []minio.ObjectInfo{
		{Err: errors.New("Error")},
		{Key: "dir2/subdir/file2.txt"},
		{Key: "dir1/file3.txt"},
	}
	suite.minioAdapter.On("ListObjects", mock.Anything, "bucketName", mock.AnythingOfType("minio.ListObjectsOptions")).Return(suite.createChan(objects))

	// when
	result, err := suite.storageService.GetBucketStructure("bucketName", "")

	// then
	suite.minioAdapter.AssertCalled(suite.T(), "ListObjects", mock.Anything, "bucketName", mock.AnythingOfType("minio.ListObjectsOptions"))
	suite.Error(err)
	suite.Nil(result)
}

func (suite *StorageServiceTest) createChan(infos []minio.ObjectInfo) <-chan minio.ObjectInfo {
	objectsCh := make(chan minio.ObjectInfo, len(infos))
	for _, obj := range infos {
		objectsCh <- obj
	}
	defer close(objectsCh)
	return objectsCh
}

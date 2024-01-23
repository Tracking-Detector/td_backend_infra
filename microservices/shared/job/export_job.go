package job

import (
	"compress/gzip"
	"context"
	"io"
	"tds/shared/models"
	"tds/shared/service"
)

type IExportJob interface {
	Execute() error
}

type InteralExportJob struct {
	exporter       *models.Exporter
	requestRepo    models.RequestRepository
	storageService service.IStorageService
}

func NewInternalExportJob(exporter *models.Exporter, requestRepo models.RequestRepository, storageService service.IStorageService) *InteralExportJob {
	return &InteralExportJob{
		exporter:       exporter,
		requestRepo:    requestRepo,
		storageService: storageService,
	}
}

func (j *InteralExportJob) Execute() error {
	ctx := context.TODO()
	pr, pw := io.Pipe()
	gzipWriter := gzip.NewWriter(pw)
	defer pw.Close()
	defer gzipWriter.Close()

}

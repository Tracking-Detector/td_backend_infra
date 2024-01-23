package job

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"

	"strings"
	"tds/shared/configs"
	"tds/shared/converter"
	"tds/shared/extractor"
	"tds/shared/models"
	"tds/shared/service"

	log "github.com/sirupsen/logrus"
)

type IExportJob interface {
	Execute() error
}

type InteralExportJob struct {
	exporter       *models.Exporter
	reducer        string
	dataset        string
	requestRepo    models.RequestRepository
	storageService service.IStorageService
}

func NewInternalExportJob(exporter *models.Exporter, reducer string, dataset string, requestRepo models.RequestRepository, storageService service.IStorageService) *InteralExportJob {
	return &InteralExportJob{
		exporter:       exporter,
		reducer:        reducer,
		dataset:        dataset,
		requestRepo:    requestRepo,
		storageService: storageService,
	}
}

func (j *InteralExportJob) Execute() error {
	ctx := context.TODO()
	extractor, err := j.getCorrectExporter(j.exporter)
	if err != nil {
		return err
	}
	pr, pw := io.Pipe()
	gzipWriter := gzip.NewWriter(pw)
	defer pw.Close()
	defer gzipWriter.Close()
	resultChannel, errorChannel := j.requestRepo.StreamByDataset(ctx, j.dataset)
	go func() {
		for {
			select {
			case requestData, ok := <-resultChannel:
				if !ok {
					break
				}
				reduced := converter.ConvertRequestModel(requestData, converter.ReduceType(j.reducer))
				encoded, encodedErr := extractor.Encode(*reduced)
				if encodedErr != nil {
					continue
				}
				arr, err := json.Marshal(encoded)
				if err != nil {
					log.WithFields(log.Fields{
						"service": "InternalExportJob",
						"error":   err.Error(),
					}).Error("Could not convert int[] to string.")
					continue
				}
				data := strings.Trim(string(arr), "[]") + "\n"
				if _, err := gzipWriter.Write([]byte(data)); err != nil {
					log.WithFields(log.Fields{
						"service": "InternalExportJob",
						"error":   err.Error(),
					}).Error("Failed to write to gzip writer.")
					break
				}
			case err, ok := <-errorChannel:
				if !ok {
					break
				}
				// Handle the error
				log.Println("Error:", err)
			}
		}
	}()
	return j.storageService.PutObject(ctx, configs.EnvExportBucketName(), j.exporter.Name+"_"+j.reducer+"_"+j.dataset+".csv.gz", pr, -1, "application/gzip")
}

func (j *InteralExportJob) getCorrectExporter(exporter *models.Exporter) (*extractor.Extractor, error) {
	for _, ext := range extractor.EXTRACTORS {
		if ext.GetName() == exporter.Name {
			return &ext, nil
		}
	}
	return nil, errors.New("In Services Extractor could not be found.")
}

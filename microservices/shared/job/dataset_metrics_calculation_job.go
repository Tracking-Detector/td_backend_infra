package job

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/service"
	log "github.com/sirupsen/logrus"
)

type IJob interface {
	Execute()
}

type DatasetMetricJob struct {
	datasetService service.IDatasetService
	requestService service.IRequestService
}

func NewDatasetMetricJob(datasetService service.IDatasetService, requestService service.IRequestService) *DatasetMetricJob {
	return &DatasetMetricJob{
		datasetService: datasetService,
		requestService: requestService,
	}
}

func (j *DatasetMetricJob) Execute() {
	ctx := context.Background()
	log.WithFields(log.Fields{
		"service": "DatasetMetricJob",
	}).Info("DatasetMetricJob started...")

	datasets := j.datasetService.GetAllDatasets()
	log.WithFields(log.Fields{
		"service": "DatasetMetricJob",
	}).Info("Loaded all datasets...")

	requestChan, errorChan := j.requestService.StreamAll(ctx)
	log.WithFields(log.Fields{
		"service": "DatasetMetricJob",
	}).Info("Loaded all requests...")
	for _, dataset := range datasets {
		dataset.Metrics = &models.DataSetMetrics{
			Total:         0,
			ReducerMetric: make([]*models.ReducerMetric, 0),
		}
	}
loop:
	for {
		select {
		case request, ok := <-requestChan:
			if !ok {
				log.WithFields(log.Fields{
					"service": "DatasetMetricJob",
				}).Info("Finished streaming requests...")
				break loop
			}
			for _, dataset := range datasets {
				if dataset.Label == request.Dataset {
					dataset.Metrics.Total++
					for _, label := range request.Labels {
						metric := j.FindMetricForLabel(dataset.Metrics.ReducerMetric, label.Blocklist)
						if metric == nil {
							metric = &models.ReducerMetric{
								Reducer: label.Blocklist,
								Total:   0,
							}
							dataset.Metrics.ReducerMetric = append(dataset.Metrics.ReducerMetric, metric)
						}
						metric.Total++
						if label.IsLabeled {
							metric.Tracker++
						} else {
							metric.NonTracker++
						}
					}
				}
			}
		case err, ok := <-errorChan:
			if !ok {
				log.WithFields(log.Fields{
					"service": "DatasetMetricJob",
				}).Info("Finished streaming requests...")
				return
			}
			log.WithFields(log.Fields{
				"service": "DatasetMetricJob",
				"error":   err,
			}).Error("Error while streaming requests...")
		}
	}
	j.datasetService.SaveAll(ctx, datasets)
	log.WithFields(log.Fields{
		"service": "DatasetMetricJob",
	}).Info("DatasetMetricJob finished...")
}

func (j *DatasetMetricJob) FindMetricForLabel(metrics []*models.ReducerMetric, label string) *models.ReducerMetric {
	for _, reducerMetric := range metrics {
		if reducerMetric.Reducer == label {
			return reducerMetric
		}
	}
	return nil
}

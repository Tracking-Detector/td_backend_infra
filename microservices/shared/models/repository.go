package models

import (
	"context"
)

type ExporterRepository interface {
	Save(ctx context.Context, p *Exporter) error
	FindByID(ctx context.Context, id string) (*Exporter, error)
	FindByName(ctx context.Context, name string) (*Exporter, error)
	FindAll(ctx context.Context) ([]*Exporter, error)
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type RequestRepository interface {
	Save(ctx context.Context, c *RequestData) error
	FindByID(ctx context.Context, id string) (*RequestData, error)
	InsertMany(ctx context.Context, requests []*RequestData) error
	CountDocuments(ctx context.Context, url string) (int64, error)
	FindAllFilteredByUrlPaged(ctx context.Context, url string, page, pageSize int) ([]*RequestData, error)
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type TrainingRunRepository interface {
	FindAll(ctx context.Context) ([]*TrainingRun, error)
	FindByModelName(ctx context.Context, modelName string) ([]*TrainingRun, error)
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

type UserRepository interface {
	Save(ctx context.Context, c *UserData) error
	DeleteUserByID(ctx context.Context, id string) error
	FindUserByID(ctx context.Context, id string) (*UserData, error)
	FindUserByEmail(ctx context.Context, email string) (*UserData, error)
	FindAll(ctx context.Context) ([]*UserData, error)
	InTransaction(ctx context.Context, fn func(context.Context) error) error
}

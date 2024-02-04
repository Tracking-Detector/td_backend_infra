package repository

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/mongostore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoExportRunRunRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoExportRunRunRepository(db *mongo.Database) *MongoExportRunRunRepository {
	coll := db.Collection(configs.EnvExporterRunsCollection())
	return &MongoExportRunRunRepository{
		db:   db,
		coll: coll,
	}
}

func (r *MongoExportRunRunRepository) Save(ctx context.Context, m *models.ExportRun) (*models.ExportRun, error) {
	return mongostore.Save(ctx, r.coll, m)
}

func (r *MongoExportRunRunRepository) SaveAll(ctx context.Context, m []*models.ExportRun) ([]*models.ExportRun, error) {
	return mongostore.SaveAll(ctx, r.coll, m)
}

func (r *MongoExportRunRunRepository) StreamAll(ctx context.Context) (<-chan *models.ExportRun, <-chan error) {
	return mongostore.StreamAll[*models.ExportRun](ctx, r.coll, bson.M{})
}

func (r *MongoExportRunRunRepository) FindByID(ctx context.Context, id string) (*models.ExportRun, error) {
	return mongostore.FindByID(ctx, r.coll, id, &models.ExportRun{})
}

func (r *MongoExportRunRunRepository) FindByExporterID(ctx context.Context, exporterId string) ([]*models.ExportRun, error) {
	return mongostore.FindAllBy(ctx, r.coll, &models.ExportRun{}, bson.M{
		"exporterId": exporterId,
	}, nil)
}

func (r *MongoExportRunRunRepository) FindAll(ctx context.Context) ([]*models.ExportRun, error) {
	return mongostore.FindAll(ctx, r.coll, &models.ExportRun{})
}

func (r *MongoExportRunRunRepository) DeleteAll(ctx context.Context) error {
	return mongostore.DeleteAll(ctx, r.coll)
}

func (r *MongoExportRunRunRepository) DeleteByID(ctx context.Context, id string) error {
	return mongostore.DeleteByID(ctx, r.coll, id)
}
func (r *MongoExportRunRunRepository) DeleteAllByExporterID(ctx context.Context, exporterId string) error {
	return mongostore.DeleteAllBy(ctx, r.coll, bson.M{
		"exporterId": exporterId,
	})
}

func (r *MongoExportRunRunRepository) Count(ctx context.Context) (int64, error) {
	return mongostore.Count(ctx, r.coll)
}

func (r *MongoExportRunRunRepository) CountByExporterID(ctx context.Context, exporterId string) (int64, error) {
	return mongostore.CountBy(ctx, r.coll, bson.M{
		"exporterId": exporterId,
	})
}

func (r *MongoExportRunRunRepository) ExistByExporterIDAndRecducer(ctx context.Context, exporterId, reducer string) (bool, error) {
	exports, err := mongostore.FindAllBy[*models.ExportRun](ctx, r.coll, &models.ExportRun{}, bson.M{
		"exporterId": exporterId,
		"reducer":    reducer,
	}, nil)
	if err != nil {
		return false, err
	}
	for _, export := range exports {
		if !export.End.IsZero() {
			return true, nil
		}
	}
	return false, nil
}

func (r *MongoExportRunRunRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return mongostore.InTransaction(ctx, r.db, fn)
}

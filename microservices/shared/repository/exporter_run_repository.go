package repository

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/mongostore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoExporterRunRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoExporterRunRepository(db *mongo.Database) *MongoExporterRunRepository {
	coll := db.Collection(configs.EnvExporterRunsCollection())
	return &MongoExporterRunRepository{
		db:   db,
		coll: coll,
	}
}

func (r *MongoExporterRunRepository) Save(ctx context.Context, m *models.Exporter) (*models.Exporter, error) {
	return mongostore.Save(ctx, r.coll, m)
}

func (r *MongoExporterRunRepository) SaveAll(ctx context.Context, m []*models.Exporter) ([]*models.Exporter, error) {
	return mongostore.SaveAll(ctx, r.coll, m)
}

func (r *MongoExporterRunRepository) StreamAll(ctx context.Context) (<-chan *models.Exporter, <-chan error) {
	return mongostore.StreamAll[*models.Exporter](ctx, r.coll, bson.M{})
}

func (r *MongoExporterRunRepository) FindByID(ctx context.Context, id string) (*models.Exporter, error) {
	return mongostore.FindByID(ctx, r.coll, id, &models.Exporter{})
}

func (r *MongoExporterRunRepository) FindByExporterID(ctx context.Context, exporterId string) ([]*models.Exporter, error) {
	return mongostore.FindAllBy(ctx, r.coll, &models.Exporter{}, bson.M{
		"exporterId": exporterId,
	}, nil)
}

func (r *MongoExporterRunRepository) FindAll(ctx context.Context) ([]*models.Exporter, error) {
	return mongostore.FindAll(ctx, r.coll, &models.Exporter{})
}

func (r *MongoExporterRunRepository) DeleteAll(ctx context.Context) error {
	return mongostore.DeleteAll(ctx, r.coll)
}

func (r *MongoExporterRunRepository) DeleteByID(ctx context.Context, id string) error {
	return mongostore.DeleteByID(ctx, r.coll, id)
}
func (r *MongoExporterRunRepository) DeleteAllByExporterID(ctx context.Context, exporterId string) error {
	return mongostore.DeleteAllBy(ctx, r.coll, bson.M{
		"exporterId": exporterId,
	})
}

func (r *MongoExporterRunRepository) Count(ctx context.Context) (int64, error) {
	return mongostore.Count(ctx, r.coll)
}

func (r *MongoExporterRunRepository) CountByExporterID(ctx context.Context, exporterId string) (int64, error) {
	return mongostore.CountBy(ctx, r.coll, bson.M{
		"exporterId": exporterId,
	})
}

func (r *MongoExporterRunRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return mongostore.InTransaction(ctx, r.db, fn)
}

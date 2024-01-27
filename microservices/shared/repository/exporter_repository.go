package repository

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/mongostore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoExporterRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoExporterRepository(db *mongo.Database) *MongoExporterRepository {
	coll := db.Collection(configs.EnvExporterCollection())
	return &MongoExporterRepository{
		db:   db,
		coll: coll,
	}
}

func (r *MongoExporterRepository) Save(ctx context.Context, m *models.Exporter) (*models.Exporter, error) {
	return mongostore.Save(ctx, r.coll, m)
}

func (r *MongoExporterRepository) SaveAll(ctx context.Context, m []*models.Exporter) ([]*models.Exporter, error) {
	return mongostore.SaveAll(ctx, r.coll, m)
}

func (r *MongoExporterRepository) StreamAll(ctx context.Context) (<-chan *models.Exporter, <-chan error) {
	return mongostore.StreamAll[*models.Exporter](ctx, r.coll, bson.M{})
}

func (r *MongoExporterRepository) FindByID(ctx context.Context, id string) (*models.Exporter, error) {
	return mongostore.FindByID(ctx, r.coll, id, &models.Exporter{})
}

func (r *MongoExporterRepository) FindAll(ctx context.Context) ([]*models.Exporter, error) {
	return mongostore.FindAll(ctx, r.coll, &models.Exporter{})
}

func (r *MongoExporterRepository) FindByName(ctx context.Context, name string) (*models.Exporter, error) {
	return mongostore.FindByName(ctx, r.coll, name, &models.Exporter{})
}

func (r *MongoExporterRepository) DeleteAll(ctx context.Context) error {
	return mongostore.DeleteAll(ctx, r.coll)
}

func (r *MongoExporterRepository) DeleteByID(ctx context.Context, id string) error {
	return mongostore.DeleteByID(ctx, r.coll, id)
}

func (r *MongoExporterRepository) Count(ctx context.Context) (int64, error) {
	return mongostore.Count(ctx, r.coll)
}

func (r *MongoExporterRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return mongostore.InTransaction(ctx, r.db, fn)
}

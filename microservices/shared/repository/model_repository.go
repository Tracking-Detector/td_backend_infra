package repository

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/mongostore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoModelRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoModelRepository(db *mongo.Database) *MongoModelRepository {
	coll := db.Collection(configs.EnvModelCollection())
	return &MongoModelRepository{
		db:   db,
		coll: coll,
	}
}

func (r *MongoModelRepository) Save(ctx context.Context, m *models.Model) (*models.Model, error) {
	return mongostore.Save(ctx, r.coll, m)
}

func (r *MongoModelRepository) SaveAll(ctx context.Context, m []*models.Model) ([]*models.Model, error) {
	return mongostore.SaveAll(ctx, r.coll, m)
}

func (r *MongoModelRepository) FindAll(ctx context.Context) ([]*models.Model, error) {
	return mongostore.FindAll(ctx, r.coll, &models.Model{})
}

func (r *MongoModelRepository) StreamAll(ctx context.Context) (<-chan *models.Model, <-chan error) {
	return mongostore.StreamAll[*models.Model](ctx, r.coll, bson.M{})
}

func (r *MongoModelRepository) FindByID(ctx context.Context, id string) (*models.Model, error) {
	return mongostore.FindByID(ctx, r.coll, id, &models.Model{})
}

func (r *MongoModelRepository) FindByName(ctx context.Context, name string) (*models.Model, error) {
	return mongostore.FindByName(ctx, r.coll, name, &models.Model{})
}

func (r *MongoModelRepository) DeleteByID(ctx context.Context, id string) error {
	return mongostore.DeleteByID(ctx, r.coll, id)
}

func (r *MongoModelRepository) DeleteAll(ctx context.Context) error {
	return mongostore.DeleteAll(ctx, r.coll)
}

func (r *MongoModelRepository) Count(ctx context.Context) (int64, error) {
	return mongostore.Count(ctx, r.coll)
}

func (r *MongoModelRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return mongostore.InTransaction(ctx, r.db, fn)
}

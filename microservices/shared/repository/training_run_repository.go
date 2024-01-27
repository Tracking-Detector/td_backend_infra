package repository

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"
	"tds/shared/mongostore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTrainingRunRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoTrainingRunRepository(db *mongo.Database) *MongoTrainingRunRepository {
	coll := db.Collection(configs.EnvTrainingRunCollection())
	mongostore.EnsureIndex(context.Background(), coll, "modelId", 1)
	mongostore.EnsureIndex(context.Background(), coll, "name", 1)
	return &MongoTrainingRunRepository{
		db:   db,
		coll: coll,
	}
}

func (r *MongoTrainingRunRepository) Save(ctx context.Context, m *models.TrainingRun) (*models.TrainingRun, error) {
	return mongostore.Save(ctx, r.coll, m)
}

func (r *MongoTrainingRunRepository) SaveAll(ctx context.Context, m []*models.TrainingRun) ([]*models.TrainingRun, error) {
	return mongostore.SaveAll(ctx, r.coll, m)
}

func (r *MongoTrainingRunRepository) FindAll(ctx context.Context) ([]*models.TrainingRun, error) {
	return mongostore.FindAll(ctx, r.coll, &models.TrainingRun{})
}

func (r *MongoTrainingRunRepository) StreamAll(ctx context.Context) (<-chan *models.TrainingRun, <-chan error) {
	return mongostore.StreamAll[*models.TrainingRun](ctx, r.coll, bson.M{})
}

func (r *MongoTrainingRunRepository) FindByID(ctx context.Context, id string) (*models.TrainingRun, error) {
	return mongostore.FindByID(ctx, r.coll, id, &models.TrainingRun{})
}

func (r *MongoTrainingRunRepository) FindByModelID(ctx context.Context, modelId string) ([]*models.TrainingRun, error) {
	return mongostore.FindAllBy(ctx, r.coll, &models.TrainingRun{}, bson.M{
		"modelId": modelId,
	}, nil)
}

func (r *MongoTrainingRunRepository) FindByName(ctx context.Context, name string) ([]*models.TrainingRun, error) {
	return mongostore.FindAllBy(ctx, r.coll, &models.TrainingRun{}, bson.M{
		"name": name,
	}, nil)
}

func (r *MongoTrainingRunRepository) DeleteByID(ctx context.Context, id string) error {
	return mongostore.DeleteByID(ctx, r.coll, id)
}

func (r *MongoTrainingRunRepository) DeleteAll(ctx context.Context) error {
	return mongostore.DeleteAll(ctx, r.coll)
}

func (r *MongoTrainingRunRepository) DeleteAllByModelID(ctx context.Context, id string) error {
	return mongostore.DeleteAllBy(ctx, r.coll, bson.M{
		"modelId": id,
	})
}

func (r *MongoTrainingRunRepository) Count(ctx context.Context) (int64, error) {
	return mongostore.Count(ctx, r.coll)
}

func (r *MongoTrainingRunRepository) CountByModelID(ctx context.Context, id string) (int64, error) {
	return mongostore.CountBy(ctx, r.coll, bson.M{
		"modelId": id,
	})
}

func (r *MongoTrainingRunRepository) CountByName(ctx context.Context, modelName string) (int64, error) {
	return mongostore.CountBy(ctx, r.coll, bson.M{
		"name": modelName,
	})
}

func (r *MongoTrainingRunRepository) FindByModelName(ctx context.Context, modelName string) ([]*models.TrainingRun, error) {
	return mongostore.FindAllBy(ctx, r.coll, &models.TrainingRun{}, bson.M{
		"name": modelName,
	}, nil)
}

func (r *MongoTrainingRunRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return mongostore.InTransaction(ctx, r.db, fn)
}

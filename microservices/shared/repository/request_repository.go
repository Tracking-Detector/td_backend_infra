package repository

import (
	"context"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/configs"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/mongostore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRequestRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func NewMongoRequestRepository(db *mongo.Database) *MongoRequestRepository {
	coll := db.Collection(configs.EnvRequestCollection())
	mongostore.EnsureIndex(context.Background(), coll, "dataset", 1)
	return &MongoRequestRepository{
		db:   db,
		coll: coll,
	}
}

func (r *MongoRequestRepository) Save(ctx context.Context, m *models.RequestData) (*models.RequestData, error) {
	return mongostore.Save(ctx, r.coll, m)
}

func (r *MongoRequestRepository) SaveAll(ctx context.Context, m []*models.RequestData) ([]*models.RequestData, error) {
	return mongostore.SaveAll(ctx, r.coll, m)
}

func (r *MongoRequestRepository) FindAll(ctx context.Context) ([]*models.RequestData, error) {
	return mongostore.FindAll(ctx, r.coll, &models.RequestData{})
}

func (r *MongoRequestRepository) FindAllByUrlLikePaged(ctx context.Context, url string, page, pageSize int) ([]*models.RequestData, error) {
	findOptions := options.Find()
	filter := bson.M{}
	if url != "" {
		filter = bson.M{
			"url": bson.M{
				"$regex": primitive.Regex{
					Pattern: url,
					Options: "i",
				},
			},
		}
	}
	findOptions.SetSkip((int64(page) - 1) * int64(pageSize))
	findOptions.SetLimit(int64(pageSize))
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))

	requestDataValues, err := mongostore.FindAllBy(ctx, r.coll, &models.RequestData{}, filter, findOptions)
	if err != nil {
		return nil, err
	}

	return requestDataValues, nil
}

func (r *MongoRequestRepository) StreamAll(ctx context.Context) (<-chan *models.RequestData, <-chan error) {
	return mongostore.StreamAll[*models.RequestData](ctx, r.coll, bson.M{})
}

func (r *MongoRequestRepository) StreamByDataset(ctx context.Context, dataset string) (<-chan *models.RequestData, <-chan error) {
	return mongostore.StreamAll[*models.RequestData](ctx, r.coll, bson.M{
		"dataset": dataset,
	})
}

func (r *MongoRequestRepository) FindByID(ctx context.Context, id string) (*models.RequestData, error) {
	return mongostore.FindByID(ctx, r.coll, id, &models.RequestData{})
}

func (r *MongoRequestRepository) Count(ctx context.Context) (int64, error) {
	return mongostore.Count(ctx, r.coll)
}

func (r *MongoRequestRepository) CountByUrlLike(ctx context.Context, url string) (int64, error) {
	filter := bson.M{}
	if url != "" {
		filter = bson.M{
			"url": bson.M{
				"$regex": primitive.Regex{
					Pattern: url,
					Options: "i",
				},
			},
		}
	}
	return mongostore.CountBy(ctx, r.coll, filter)
}

func (r *MongoRequestRepository) DeleteByID(ctx context.Context, id string) error {
	return mongostore.DeleteByID(ctx, r.coll, id)
}

func (r *MongoRequestRepository) DeleteAll(ctx context.Context) error {
	return mongostore.DeleteAll(ctx, r.coll)
}

func (r *MongoRequestRepository) DeleteAllByLabel(ctx context.Context, label string) error {
	return mongostore.DeleteAllBy(ctx, r.coll, bson.M{
		"dataset": label,
	})
}

func (r *MongoRequestRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return mongostore.InTransaction(ctx, r.db, fn)
}

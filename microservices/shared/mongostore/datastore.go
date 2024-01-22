package mongostore

import (
	"context"
	"tds/shared/configs"
	"tds/shared/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindByID(ctx context.Context, coll *mongo.Collection, id string, m interface{}) error {
	objId, _ := primitive.ObjectIDFromHex(id)
	res := coll.FindOne(ctx, bson.M{"_id": objId})
	if err := res.Err(); err != nil {
		return err
	}
	return res.Decode(m)
}

func CountDocuments(ctx context.Context, coll *mongo.Collection, filter bson.M) (int64, error) {
	return coll.CountDocuments(ctx, filter)
}

func InsertMany(ctx context.Context, coll *mongo.Collection, m []interface{}) error {
	_, err := coll.InsertMany(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func FindAll(ctx context.Context, coll *mongo.Collection, filter bson.M, options *options.FindOptions, results interface{}) error {
	cursor, err := coll.Find(ctx, filter, options)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

func FindExporterByID(ctx context.Context, db *mongo.Database, id string) (*models.Exporter, error) {
	p := new(models.Exporter)
	err := FindByID(ctx, db.Collection(configs.EnvExporterCollection()), id, p)
	return p, err
}

func FindAllExporter(ctx context.Context, db *mongo.Database) ([]*models.Exporter, error) {
	var exporter []*models.Exporter
	err := FindAll(ctx, db.Collection(configs.EnvExporterCollection()), bson.M{}, nil, &exporter)
	return exporter, err
}

func SaveExporter(ctx context.Context, db *mongo.Database, p *models.Exporter) error {
	opts := options.FindOneAndReplace().SetUpsert(true)
	var doc bson.M
	err := db.Collection(configs.EnvExporterCollection()).FindOneAndReplace(ctx, bson.D{{Key: "name", Value: p.Name}}, p, opts).Decode(&doc)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func SaveRequest(ctx context.Context, db *mongo.Database, p *models.RequestData) error {
	opts := options.FindOneAndReplace().SetUpsert(true)
	var doc bson.M
	err := db.Collection(configs.EnvRequestCollection()).FindOneAndReplace(ctx, bson.D{{Key: "_id", Value: p.Id}}, p, opts).Decode(&doc)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func FindRequestByID(ctx context.Context, db *mongo.Database, id string) (*models.RequestData, error) {
	p := new(models.RequestData)
	err := FindByID(ctx, db.Collection(configs.EnvRequestCollection()), id, p)
	return p, err
}

func InsertManyRequests(ctx context.Context, db *mongo.Database, requests []*models.RequestData) error {
	interfaceSlice := make([]interface{}, len(requests))
	for i, v := range requests {
		interfaceSlice[i] = v
	}
	return InsertMany(ctx, db.Collection(configs.EnvRequestCollection()), interfaceSlice)
}

func CountRequestDocumentsForUrlFilter(ctx context.Context, db *mongo.Database, url string) (int64, error) {
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
	return CountDocuments(ctx, db.Collection(configs.EnvRequestCollection()), filter)
}

func FindAllRequestFilteredByUrlPaged(ctx context.Context, db *mongo.Database, url string, page, pageSize int) ([]*models.RequestData, error) {
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

	var requestDataValues []*models.RequestData
	err := FindAll(ctx, db.Collection(configs.EnvRequestCollection()), filter, findOptions, &requestDataValues)
	if err != nil {
		return nil, err
	}

	return requestDataValues, nil
}

func FindAllTrainingRuns(ctx context.Context, db *mongo.Database) ([]*models.TrainingRun, error) {
	var trainingRuns []*models.TrainingRun
	err := FindAll(ctx, db.Collection(configs.EnvTrainingRunCollection()), bson.M{}, nil, trainingRuns)
	return trainingRuns, err
}

func FindAllTrainingRunsByModelName(ctx context.Context, db *mongo.Database, modelName string) ([]*models.TrainingRun, error) {
	var trainingRuns []*models.TrainingRun
	err := FindAll(ctx, db.Collection(configs.EnvTrainingRunCollection()), bson.M{
		"name": modelName,
	}, nil, trainingRuns)
	return trainingRuns, err
}

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
	res := coll.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		return err
	}
	return res.Decode(m)
}

func DeleteAll(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.DeleteMany(ctx, bson.M{})
	return err
}

func FindByName(ctx context.Context, coll *mongo.Collection, name string, m interface{}) error {
	res := coll.FindOne(ctx, bson.M{"name": name})
	if err := res.Err(); err != nil {
		return err
	}
	return res.Decode(m)
}

func DeleteByID(ctx context.Context, coll *mongo.Collection, id string) error {
	_, err := coll.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	return err
}

func DeleteAllBy(ctx context.Context, coll *mongo.Collection, key string, value string) error {
	_, err := coll.DeleteMany(ctx, bson.M{
		key: value,
	})
	return err
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

func DeleteAllExporter(ctx context.Context, db *mongo.Database) error {
	return DeleteAll(ctx, db.Collection(configs.EnvExporterCollection()))
}

func FindAllExporter(ctx context.Context, db *mongo.Database) ([]*models.Exporter, error) {
	var exporter []*models.Exporter
	err := FindAll(ctx, db.Collection(configs.EnvExporterCollection()), bson.M{}, nil, &exporter)
	return exporter, err
}

func FindExporterByName(ctx context.Context, db *mongo.Database, name string) (*models.Exporter, error) {
	exporter := new(models.Exporter)
	err := FindByName(ctx, db.Collection(configs.EnvExporterCollection()), name, exporter)
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
	err := db.Collection(configs.EnvRequestCollection()).FindOneAndReplace(ctx, bson.D{{Key: "_id", Value: p.ID}}, p, opts).Decode(&doc)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func DeleteAllRequests(ctx context.Context, db *mongo.Database) error {
	return DeleteAll(ctx, db.Collection(configs.EnvRequestCollection()))
}

func FindRequestByID(ctx context.Context, db *mongo.Database, id string) (*models.RequestData, error) {
	p := new(models.RequestData)
	err := FindByID(ctx, db.Collection(configs.EnvRequestCollection()), id, p)
	return p, err
}

func StreamRequestsByDataset(ctx context.Context, db *mongo.Database, dataset string) (<-chan *models.RequestData, <-chan error) {
	resultChannel := make(chan *models.RequestData)
	errorChannel := make(chan error)
	go func() {
		defer close(resultChannel)
		defer close(errorChannel)
		filter := bson.M{}
		if dataset != "" {
			filter = bson.M{
				"dataset": dataset,
			}
		}
		// Set up the cursor to traverse the entire collection
		opts := options.Find().SetCursorType(options.TailableAwait)
		cursor, err := db.Collection(configs.EnvRequestCollection()).Find(ctx, filter, opts)
		if err != nil {
			errorChannel <- err
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var requestData models.RequestData
			if err := cursor.Decode(&requestData); err != nil {
				errorChannel <- err
				return
			}

			resultChannel <- &requestData
		}

		if err := cursor.Err(); err != nil {
			errorChannel <- err
		}
	}()

	return resultChannel, errorChannel
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

func DeleteAllTrainingRuns(ctx context.Context, db *mongo.Database) error {
	return DeleteAll(ctx, db.Collection(configs.EnvTrainingRunCollection()))
}

func FindAllTrainingRunsByModelName(ctx context.Context, db *mongo.Database, modelName string) ([]*models.TrainingRun, error) {
	var trainingRuns []*models.TrainingRun
	err := FindAll(ctx, db.Collection(configs.EnvTrainingRunCollection()), bson.M{
		"name": modelName,
	}, nil, trainingRuns)
	return trainingRuns, err
}

func FindAllByModelId(ctx context.Context, db *mongo.Database, modelId string) ([]*models.TrainingRun, error) {
	var trainingRuns []*models.TrainingRun
	err := FindAll(ctx, db.Collection(configs.EnvTrainingRunCollection()), bson.M{
		"modelId": modelId,
	}, nil, trainingRuns)
	return trainingRuns, err
}

func DeleteTrainingRunById(ctx context.Context, db *mongo.Database, id string) error {
	return DeleteByID(ctx, db.Collection(configs.EnvTrainingRunCollection()), id)
}

func DeleteTrainingRunsByModelId(ctx context.Context, db *mongo.Database, modelId string) error {
	return DeleteAllBy(ctx, db.Collection(configs.EnvTrainingRunCollection()), "modelId", modelId)
}

func SaveUser(ctx context.Context, db *mongo.Database, p *models.UserData) error {
	opts := options.FindOneAndReplace().SetUpsert(true)
	var doc bson.M
	err := db.Collection(configs.EnvUserCollection()).FindOneAndReplace(ctx, bson.D{{Key: "email", Value: p.Email}}, p, opts).Decode(&doc)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func DeleteUserByID(ctx context.Context, db *mongo.Database, id string) error {
	objId, _ := primitive.ObjectIDFromHex(id)
	_, err := db.Collection(configs.EnvUserCollection()).DeleteOne(ctx, bson.M{
		"_id": objId,
	}, nil)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func FindAllUsers(ctx context.Context, db *mongo.Database) ([]*models.UserData, error) {
	var users []*models.UserData
	err := FindAll(ctx, db.Collection(configs.EnvUserCollection()), bson.M{}, nil, users)
	return users, err
}

func DeleteAllUser(ctx context.Context, db *mongo.Database) error {
	return DeleteAll(ctx, db.Collection(configs.EnvUserCollection()))
}

func FindUserByID(ctx context.Context, db *mongo.Database, id string) (*models.UserData, error) {
	p := new(models.UserData)
	err := FindByID(ctx, db.Collection(configs.EnvUserCollection()), id, p)
	return p, err
}

func FindUserByEmail(ctx context.Context, db *mongo.Database, email string) (*models.UserData, error) {
	user := new(models.UserData)
	err := db.Collection(configs.EnvUserCollection()).FindOne(ctx, bson.M{
		"email": email,
	}).Decode(user)
	return user, err
}

func SaveModel(ctx context.Context, db *mongo.Database, model *models.Model) error {
	opts := options.FindOneAndReplace().SetUpsert(true)
	var doc bson.M
	err := db.Collection(configs.EnvModelCollection()).FindOneAndReplace(ctx, bson.D{{Key: "name", Value: model.Name}}, model, opts).Decode(&doc)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func DeleteAllModels(ctx context.Context, db *mongo.Database) error {
	return DeleteAll(ctx, db.Collection(configs.EnvModelCollection()))
}

func DeleteModelByID(ctx context.Context, db *mongo.Database, id string) error {
	return DeleteByID(ctx, db.Collection(configs.EnvModelCollection()), id)
}

func FindModelByID(ctx context.Context, db *mongo.Database, id string) (*models.Model, error) {
	p := new(models.Model)
	err := FindByID(ctx, db.Collection(configs.EnvModelCollection()), id, p)
	return p, err
}

func FindModelByName(ctx context.Context, db *mongo.Database, name string) (*models.Model, error) {
	p := new(models.Model)
	err := FindByName(ctx, db.Collection(configs.EnvModelCollection()), name, p)
	return p, err
}

func FindAllModels(ctx context.Context, db *mongo.Database) ([]*models.Model, error) {
	var models []*models.Model
	err := FindAll(ctx, db.Collection(configs.EnvUserCollection()), bson.M{}, nil, models)
	return models, err
}

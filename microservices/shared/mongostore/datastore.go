package mongostore

import (
	"context"
	"fmt"

	"github.com/Tracking-Detector/td_backend_infra/microservices/shared/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndex(ctx context.Context, coll *mongo.Collection, name string, direction int) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: name, Value: direction},
		},
		Options: options.Index().SetName(fmt.Sprintf("%s_index", name)),
	}

	_, err := coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}

	return nil
}

// Base Interface Functions

func Save[T models.BaseModel](ctx context.Context, coll *mongo.Collection, entity T) (T, error) {
	if entity.GetID() == "" {
		entity.SetID(uuid.New().String())
	}
	opts := options.FindOneAndReplace().SetUpsert(true)

	err := coll.FindOneAndReplace(ctx, bson.D{{Key: "_id", Value: entity.GetID()}}, entity, opts).Decode(entity)
	if err != nil && err != mongo.ErrNoDocuments {
		return entity, err
	}

	return entity, nil
}

func Exist(ctx context.Context, coll *mongo.Collection, filter bson.M) (bool, error) {
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func SaveAll[T models.BaseModel](ctx context.Context, coll *mongo.Collection, entities []T) ([]T, error) {
	var savedEntities []T

	// Prepare the bulk write operations
	var bulkWrites []mongo.WriteModel
	for _, entity := range entities {
		// Generate a random UUID if _id is not set
		if entity.GetID() == "" {
			entity.SetID(uuid.New().String())
		}

		filter := bson.D{{Key: "_id", Value: entity.GetID()}}
		upsert := true
		bulkWrite := mongo.NewReplaceOneModel().SetFilter(filter).SetReplacement(entity).SetUpsert(upsert)
		bulkWrites = append(bulkWrites, bulkWrite)
	}

	// Execute the bulk write operations
	opts := options.BulkWrite().SetOrdered(false) // SetOrdered(false) for unordered execution
	result, err := coll.BulkWrite(ctx, bulkWrites, opts)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated entities after the bulk write
	for _, upsertedID := range result.UpsertedIDs {
		// Ensure that the upsertedID is a valid ObjectID
		id, ok := upsertedID.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for upsertedID")
		}

		// Find the corresponding entity using the ObjectID
		for _, entity := range entities {
			if entity.GetID() == id {
				savedEntities = append(savedEntities, entity)
				break
			}
		}
	}

	return savedEntities, nil
}

func FindAll[T models.BaseModel](ctx context.Context, coll *mongo.Collection, entityType T) ([]T, error) {
	var results []T
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FindAllBy[T models.BaseModel](ctx context.Context, coll *mongo.Collection, entityType T, filter bson.M, findOptions *options.FindOptions) ([]T, error) {
	var results []T
	cursor, err := coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FindByID[T models.BaseModel](ctx context.Context, coll *mongo.Collection, id string, entityType T) (T, error) {
	var entity T
	res := coll.FindOne(ctx, bson.M{"_id": id})
	if err := res.Err(); err != nil {
		return entity, err
	}
	if err := res.Decode(entityType); err != nil {
		return entity, err
	}

	return entityType, nil
}

func FindBy[T models.BaseModel](ctx context.Context, coll *mongo.Collection, filter bson.M, entityType T) (T, error) {
	var entity T
	res := coll.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return entity, err
	}
	if err := res.Decode(entityType); err != nil {
		return entity, err
	}

	return entityType, nil
}

func FindByName[T models.BaseModelName](ctx context.Context, coll *mongo.Collection, name string, entityType T) (T, error) {
	var entity T
	res := coll.FindOne(ctx, bson.M{"name": name})
	if err := res.Err(); err != nil {
		return entity, err
	}

	if err := res.Decode(entityType); err != nil {
		return entity, err
	}

	return entityType, nil
}

func DeleteByID(ctx context.Context, coll *mongo.Collection, id string) error {
	_, err := coll.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	return err
}

func DeleteAll(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.DeleteMany(ctx, bson.M{})
	return err
}

func Count(ctx context.Context, coll *mongo.Collection) (int64, error) {
	return coll.CountDocuments(ctx, bson.M{})
}

func CountBy(ctx context.Context, coll *mongo.Collection, filter bson.M) (int64, error) {
	return coll.CountDocuments(ctx, filter)
}

func StreamAll[T models.BaseModel](ctx context.Context, db *mongo.Collection, filter bson.M) (<-chan T, <-chan error) {
	resultChannel := make(chan T)
	errorChannel := make(chan error)

	go func() {
		defer close(resultChannel)
		defer close(errorChannel)

		opts := options.Find()
		cursor, err := db.Find(ctx, filter, opts)
		if err != nil {
			fmt.Println(err)
			errorChannel <- err
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var data T
			if err := cursor.Decode(&data); err != nil {
				fmt.Println(err)
				errorChannel <- err
				return
			}

			resultChannel <- data
		}

		if err := cursor.Err(); err != nil {
			fmt.Println(err)
			errorChannel <- err
		}
	}()

	return resultChannel, errorChannel
}

func InTransaction(ctx context.Context, db *mongo.Database, fn func(context.Context) error) error {
	sess, err := db.Client().StartSession()
	if err != nil {
		return err
	}

	return mongo.WithSession(ctx, sess, func(sc mongo.SessionContext) error {
		defer sess.EndSession(context.Background())

		if err := sc.StartTransaction(); err != nil {
			return err
		}
		if err := fn(sc); err != nil {
			return sc.AbortTransaction(sc)
		}
		return sc.CommitTransaction(sc)
	})
}

func DeleteAllBy(ctx context.Context, coll *mongo.Collection, filter bson.M) error {
	_, err := coll.DeleteMany(ctx, filter)
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

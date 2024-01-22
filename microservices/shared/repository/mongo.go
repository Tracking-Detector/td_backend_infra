package repository

import (
	"context"
	"tds/shared/models"
	"tds/shared/mongostore"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoExporterRepository struct {
	db *mongo.Database
}

func inMongoTransaction(ctx context.Context, db *mongo.Database, fn func(context.Context) error) error {
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

func NewMongoExporterRepository(db *mongo.Database) *MongoExporterRepository {
	return &MongoExporterRepository{db: db}
}

func (r *MongoExporterRepository) FindByID(ctx context.Context, id string) (*models.Exporter, error) {
	return mongostore.FindExporterByID(ctx, r.db, id)
}

func (r *MongoExporterRepository) Save(ctx context.Context, m *models.Exporter) error {
	return mongostore.SaveExporter(ctx, r.db, m)
}

func (r *MongoExporterRepository) FindAll(ctx context.Context) ([]*models.Exporter, error) {
	return mongostore.FindAllExporter(ctx, r.db)
}

func (r *MongoExporterRepository) FindByName(ctx context.Context, name string) (*models.Exporter, error) {
	return mongostore.FindExporterByName(ctx, r.db, name)
}

func (r *MongoExporterRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return inMongoTransaction(ctx, r.db, fn)
}

type MongoRequestRepository struct {
	db *mongo.Database
}

func NewMongoRequestRepository(db *mongo.Database) *MongoRequestRepository {
	return &MongoRequestRepository{
		db: db,
	}
}

func (r *MongoRequestRepository) FindByID(ctx context.Context, id string) (*models.RequestData, error) {
	return mongostore.FindRequestByID(ctx, r.db, id)
}

func (r *MongoRequestRepository) InsertMany(ctx context.Context, requests []*models.RequestData) error {
	return mongostore.InsertManyRequests(ctx, r.db, requests)
}

func (r *MongoRequestRepository) Save(ctx context.Context, m *models.RequestData) error {
	return mongostore.SaveRequest(ctx, r.db, m)
}

func (r *MongoRequestRepository) CountDocuments(ctx context.Context, url string) (int64, error) {
	return mongostore.CountRequestDocumentsForUrlFilter(ctx, r.db, url)
}

func (r *MongoRequestRepository) FindAllFilteredByUrlPaged(ctx context.Context, url string, page, pageSize int) ([]*models.RequestData, error) {
	return mongostore.FindAllRequestFilteredByUrlPaged(ctx, r.db, url, page, pageSize)
}

func (r *MongoRequestRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return inMongoTransaction(ctx, r.db, fn)
}

type MongoTrainingRunsRepository struct {
	db *mongo.Database
}

func NewMongoTrainingRunsRepository(db *mongo.Database) *MongoTrainingRunsRepository {
	return &MongoTrainingRunsRepository{
		db: db,
	}
}

func (r *MongoTrainingRunsRepository) FindAll(ctx context.Context) ([]*models.TrainingRun, error) {
	return mongostore.FindAllTrainingRuns(ctx, r.db)
}

func (r *MongoTrainingRunsRepository) FindByModelName(ctx context.Context, modelName string) ([]*models.TrainingRun, error) {
	return mongostore.FindAllTrainingRunsByModelName(ctx, r.db, modelName)
}

func (r *MongoTrainingRunsRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return inMongoTransaction(ctx, r.db, fn)
}

type MongoUserRepository struct {
	db *mongo.Database
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		db: db,
	}
}

func (r *MongoUserRepository) Save(ctx context.Context, m *models.UserData) error {
	return mongostore.SaveUser(ctx, r.db, m)
}

func (r *MongoUserRepository) DeleteUserByID(ctx context.Context, id string) error {
	return mongostore.DeleteUserByID(ctx, r.db, id)
}

func (r *MongoUserRepository) FindAll(ctx context.Context) ([]*models.UserData, error) {
	return mongostore.FindAllUsers(ctx, r.db)
}

func (r *MongoUserRepository) FindUserByID(ctx context.Context, id string) (*models.UserData, error) {
	return mongostore.FindUserByID(ctx, r.db, id)
}

func (r *MongoUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.UserData, error) {
	return mongostore.FindUserByEmail(ctx, r.db, email)
}

func (r *MongoUserRepository) InTransaction(ctx context.Context, fn func(context.Context) error) error {
	return inMongoTransaction(ctx, r.db, fn)
}

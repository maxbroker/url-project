package storage

import (
	"awesomeProject/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type Storage struct {
	db  *mongo.Collection
	ctx context.Context
}

type StorageFile struct {
	Alias string
	URL   string
}

var (
	ErrURLExists   = errors.New("URL already exists")
	ErrURLNotFound = errors.New("URL not found")
	ZeroID         primitive.ObjectID
)

func ConnectToDB(CollectionName string, DBName string, cfg *config.Config, logger *slog.Logger, ctx context.Context) (*Storage, error) {
	const op = "migrations.ConnectToDB"
	dsn := fmt.Sprintf("mongodb://%s:%s", cfg.Dbhost, cfg.Dbport)
	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	db := client.Database(DBName)
	if err := initCollections(db, CollectionName); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	collection := db.Collection(CollectionName)
	createStorage := Storage{collection, ctx}
	logger.Info("Connected to MongoDB!")
	return &createStorage, nil
}

func RunMigrations(dbName string, cfg *config.Config) error {
	// Путь к папке с миграциями
	migrationsPath := "file://db/migrations"
	m, err := migrate.New(
		migrationsPath,
		fmt.Sprintf("mongodb://%s:%s/%s", cfg.Dbhost, cfg.Dbport, dbName),
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func initCollections(db *mongo.Database, dbName string) error {
	collections := []string{
		dbName,
	}
	for _, collName := range collections {
		if err := createCollectionIfNotExists(db, collName); err != nil {
			return err
		}
	}
	return nil
}

func createCollectionIfNotExists(db *mongo.Database, collName string) error {
	err := db.CreateCollection(context.TODO(), collName)
	if err != nil && !isCollectionExistsError(err) {
		return err
	}
	return nil
}

func isCollectionExistsError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "collection already exists"
}

func (s *Storage) SaveURL(urlToSave string, alias string) (primitive.ObjectID, error) {
	const op = "migrations.saveUrl"
	storage := s.db
	linkToSave := StorageFile{Alias: alias, URL: urlToSave}
	filter := bson.M{"alias": alias}
	var existUrl StorageFile

	err := storage.FindOne(s.ctx, filter).Decode(&existUrl)
	if err == nil {
		return ZeroID, ErrURLExists
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return ZeroID, fmt.Errorf("%s: %v", op, err)
	}
	insertResult, err := storage.InsertOne(s.ctx, linkToSave)
	if err != nil {
		return ZeroID, fmt.Errorf("%s: %v", op, err)
	}
	objectID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return ZeroID, fmt.Errorf("%s: failed to convert InsertedID to ObjectID", op)
	}
	fmt.Printf("Url was been saved", slog.Any("_id", objectID))
	return objectID, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "migrations.getUrl"
	storage := s.db
	filter := bson.M{"alias": alias}
	var existUrl StorageFile

	err := storage.FindOne(s.ctx, filter).Decode(&existUrl)
	if err != nil {
		return "", ErrURLNotFound
	}

	return existUrl.URL, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "migrations.deleteUrl"
	storage := s.db
	filter := bson.M{"alias": alias}
	result, err := storage.DeleteOne(s.ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.DeletedCount == 0 {
		return ErrURLNotFound
	}
	return nil
}

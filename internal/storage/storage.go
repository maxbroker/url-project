package storage

import (
	"awesomeProject/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb" // Импорт драйвера MongoDB
	_ "github.com/golang-migrate/migrate/v4/source/file"      // Импорт драйвера файловой системы
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

type File struct {
	Alias string
	URL   string
}

var (
	ErrURLExists   = errors.New("URL already exists")
	ErrURLNotFound = errors.New("URL not found")
	ZeroID         primitive.ObjectID
)

func ConnectToDB(collectionName string, dbName string, cfg *config.Config, logger *slog.Logger, ctx context.Context) (*Storage, error) {
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
	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	createStorage := Storage{collection, ctx}
	logger.Info("Connected to MongoDB!")
	return &createStorage, nil
}

func RunMigrations(dbName string, cfg *config.Config) error {
	migrationsPath := "file://./db/migrations"
	var connectionString string
	if cfg.Env == "local" {
		connectionString = fmt.Sprintf(
			"mongodb://%s:%s/%s",
			cfg.Dbhost,
			cfg.Dbport,
			dbName)
	} else {
		connectionString = fmt.Sprintf(
			"mongodb://%s:%s@%s:%s/%s?authSource=admin&authMechanism=SCRAM-SHA-1",
			cfg.UserDB,
			cfg.PasswordDB,
			cfg.Dbhost,
			cfg.Dbport,
			dbName,
		)
	}

	m, err := migrate.New(migrationsPath, connectionString)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (primitive.ObjectID, error) {
	const op = "migrations.saveUrl"
	storage := s.db
	linkToSave := File{Alias: alias, URL: urlToSave}
	filter := bson.M{"alias": alias}
	var existUrl File

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
	return objectID, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	storage := s.db
	filter := bson.M{"alias": alias}
	var existUrl File

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

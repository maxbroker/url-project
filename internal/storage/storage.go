package storage

import (
	"awesomeProject/internal/config"
	"context"
	"errors"
	"fmt"
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

func ConnectToDB(CollectionName string, cfg *config.Config, logger *slog.Logger, ctx context.Context) (*Storage, error) {
	const op = "storage.ConnectToDB"
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

	logger.Info("Connected to MongoDB!")
	collection := client.Database("mydb").Collection(CollectionName)
	createStorage := Storage{collection, ctx}
	return &createStorage, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (primitive.ObjectID, error) {
	const op = "storage.saveUrl"
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

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.getUrl"
	storage := s.db
	filter := bson.M{"alias": alias}
	var existUrl StorageFile

	err := storage.FindOne(s.ctx, filter).Decode(&existUrl)
	if err != nil {
		return existUrl.URL, ErrURLNotFound
	}
	return existUrl.URL, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.getUrl"
	storage := s.db
	filter := bson.M{"alias": alias}

	_, err := storage.DeleteOne(s.ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

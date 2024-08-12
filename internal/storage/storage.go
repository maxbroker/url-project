package storage

import (
	"awesomeProject/internal/config"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
)

type Storage struct {
	db *mongo.Collection
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

type StorageFile struct {
	ID    int64
	Alias string
	URL   string
}

var ErrURLExists = errors.New("URL already exists")

func ConnectingToDB(CollectionName string, cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	const op = "storage.ConnectingToDB"
	dsn := fmt.Sprintf("mongodb://%s:%s", cfg.Dbhost, cfg.Dbport)
	clientOptions := options.Client().ApplyURI(dsn)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	logger.Info("Connected to MongoDB!", slog.String("env", cfg.Env))
	collection := client.Database("mydb").Collection(CollectionName)
	createStorage := Storage{collection}
	return &createStorage, nil
}

func (s *Storage) SavingUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.saveUrl"
	storage := s.db
	count, err := storage.CountDocuments(context.TODO(), bson.M{}) // Счетчик документов в коллекции, передаем базовый контекст (в более сложных случаях можно передать какой-нибудь
	// другой контекст, который будет управлять временем и отменой операции, а также устанавливается filter, в нашем случае пустой, т.к. нужна вся коллекция)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	linkToSave := StorageFile{ID: int64(count) + 1, Alias: alias, URL: urlToSave}
	filter := bson.M{"alias": alias}
	var existUrl StorageFile

	err = storage.FindOne(context.TODO(), filter).Decode(&existUrl)
	if err == nil {
		return linkToSave.ID, ErrURLExists
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return linkToSave.ID, fmt.Errorf("%s: %v", op, err)
	}

	result, err := storage.InsertOne(context.TODO(), linkToSave)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	_ = result // говорим Go, что эта штука нам ещё понадобиться(дай бог)
	return linkToSave.ID, nil
}

func (s *Storage) GettingUrl(alias string) (string, error) {
	const op = "storage.getUrl"
	storage := s.db
	filter := bson.M{"alias": alias}
	var existUrl StorageFile

	err := storage.FindOne(context.TODO(), filter).Decode(&existUrl)
	if err != nil {
		return existUrl.URL, fmt.Errorf("%s: %v", op, err)
	}
	return existUrl.URL, nil
}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "storage.getUrl"
	storage := s.db
	filter := bson.M{"alias": alias}

	_, err := storage.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

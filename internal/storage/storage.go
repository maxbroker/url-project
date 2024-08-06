package storage

import (
	"awesomeProject/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Collection
}

type StorageFile struct {
	ID    int
	Alias string
	URL   string
}

func ConnectionToDB(CollectionName string, cfg *config.Config) (*Storage, error) {

	const op = "storage.mongoDB.ConnectionToDB"

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

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("mydb").Collection(CollectionName)

	CreateCollection := Storage{collection}
	return &CreateCollection, nil
}

func (s *Storage) SaveUrl(UrlToSave string, alias string) (int64, error) {
	const op = "storage.saveUrl"
	collection := s.db
	count, err := collection.CountDocuments(context.TODO(), bson.M{}) // Счетчик документов в коллекции, передаем базовый контекст (в более сложных случаях можно передать какой-нибудь
	// другой контекст, который будет управлять временем и отменой операции, а также устанавливается filter, в нашем случае пустой, т.к. нужна вся коллекция)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	linkToSave := StorageFile{ID: int(count) + 1, Alias: alias, URL: UrlToSave}

	InsertOneResult, err := collection.InsertOne(context.TODO(), linkToSave)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	_ = InsertOneResult // говорим Go, что эта штука нам ещё понадобиться(дай бог)
	return int64(linkToSave.ID), nil
}

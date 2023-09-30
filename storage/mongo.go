package storage

import (
	"context"
	"os"
	"time"

	"github.com/ankeshnirala/go/go_services/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	db *mongo.Database
}

func NewMongoStore() (*MongoStore, error) {
	connStr := os.Getenv("MONGODB_URL")
	DATABASE := os.Getenv("DATABASE_NAME")

	clientOptions := options.Client().ApplyURI(connStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MongoStore{db: client.Database(DATABASE)}, nil
}

func (s *MongoStore) InsertOne(collection string, user *types.User) (*mongo.InsertOneResult, error) {
	return s.db.Collection(collection).InsertOne(context.TODO(), user)
}

func (s *MongoStore) UpdateOne(collection string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return s.db.Collection(collection).UpdateOne(context.TODO(), filter, update)
}

func (s *MongoStore) DeleteByID(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	return s.db.Collection(collection).DeleteOne(context.TODO(), filter)
}

func (s *MongoStore) FindOne(collection string, filter interface{}) *mongo.SingleResult {
	return s.db.Collection(collection).FindOne(context.TODO(), filter)
}

func (s *MongoStore) Find(collection string, filter interface{}) (*mongo.Cursor, error) {
	return s.db.Collection(collection).Find(context.TODO(), filter)
}

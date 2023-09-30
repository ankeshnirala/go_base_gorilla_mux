package storage

import (
	"github.com/ankeshnirala/go/go_services/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage interface {
	InsertOne(string, *types.User) (*mongo.InsertOneResult, error)
	DeleteByID(string, interface{}) (*mongo.DeleteResult, error)
	UpdateOne(string, interface{}, interface{}) (*mongo.UpdateResult, error)
	FindOne(string, interface{}) *mongo.SingleResult
	Find(string, interface{}) (*mongo.Cursor, error)
}

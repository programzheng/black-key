package database

import (
	"context"

	"github.com/programzheng/black-key/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoInstance() *MongoInstance {
	uri := config.Cfg.GetString("MONGO_URI")
	ctx := context.TODO()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	db := c.Database(config.Cfg.GetString("MONGO_DATABASE"))
	mi := &MongoInstance{
		Client:   c,
		Database: db,
	}
	return mi
}

func (mi *MongoInstance) CreateOne(c string, m interface{}) (*mongo.InsertOneResult, error) {
	ctx := context.TODO()
	result, err := mi.Database.Collection(c).InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}
	return result, nil
}

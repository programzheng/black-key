package database

import (
	"context"

	"github.com/programzheng/black-key/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type MongoBaseRepository struct {
	MongoInstance
	CollectionName string
}

func newMongoInstance() *MongoInstance {
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

func NewMongoBaseRepository() *MongoBaseRepository {
	return &MongoBaseRepository{
		MongoInstance:  *newMongoInstance(),
		CollectionName: "",
	}
}

func (mbr *MongoBaseRepository) CreateOne(m interface{}) (*string, error) {
	ctx := context.TODO()
	r, err := mbr.MongoInstance.Database.Collection(mbr.CollectionName).InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}
	id := r.InsertedID.(primitive.ObjectID).Hex()
	return &id, nil
}

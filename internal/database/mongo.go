package database

import (
	"context"

	"github.com/programzheng/black-key/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

type MongoBaseRepository struct {
	MongoInstance
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

func NewMongoBaseRepository(collectionName string) *MongoBaseRepository {
	mi := *newMongoInstance()
	mi.Collection = mi.Database.Collection(collectionName)
	return &MongoBaseRepository{
		MongoInstance: mi,
	}
}

func (mbr *MongoBaseRepository) CreateOne(m interface{}) (*string, error) {
	ctx := context.TODO()
	r, err := mbr.MongoInstance.Collection.InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}
	id := r.InsertedID.(primitive.ObjectID).Hex()
	return &id, nil
}

func (mbr *MongoBaseRepository) Find(f interface{}, ms interface{}) (interface{}, error) {
	ctx := context.TODO()
	cursor, err := mbr.MongoInstance.Collection.Find(ctx, f)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &ms); err != nil {
		return nil, err
	}
	return ms, nil
}

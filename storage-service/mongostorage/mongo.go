package mongostorage

import (
	"context"
	"fmt"
	"time"

	"github.com/matrosov-nikita/newsapp/storage-service"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/matrosov-nikita/newsapp/storage-service/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewsMongoDocument represents how news entity stored in mongo collection.
type NewsMongoDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	Header    string             `bson:"header"`
	CreatedAt time.Time          `bson:"created_at"`
}

// NewsDocument creates new mongo document for news.
func NewsDocument(news *pb.News) *NewsMongoDocument {
	t, _ := ptypes.Timestamp(news.CreatedAt)
	return &NewsMongoDocument{ID: primitive.NewObjectID(), Header: news.Header, CreatedAt: t}
}

// ToNews converts mongo document to news.
func (d *NewsMongoDocument) ToNews() *pb.News {
	createdAt, _ := ptypes.TimestampProto(d.CreatedAt)
	return &pb.News{
		Id:        d.ID.Hex(),
		Header:    d.Header,
		CreatedAt: createdAt,
	}
}

// NewsRepository represents mongo repo for news.
type NewsRepository struct {
	c *mongo.Collection
}

// CreateNewsRepository creates new mongo repo for news.
func CreateNewsRepository(client *mongo.Client) *NewsRepository {
	return &NewsRepository{
		c: client.Database("newsdb").Collection("news"),
	}
}

// Save inserts news and return inserted id.
func (repo *NewsRepository) Save(news *pb.News) (string, error) {
	doc := NewsDocument(news)
	_, err := repo.c.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}
	return doc.ID.Hex(), nil
}

// FindByID finds news by id and returns error if not found.
func (repo *NewsRepository) FindByID(id string) (*pb.News, error) {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("incorrect id: %s", id)
	}

	var doc NewsMongoDocument
	if err := repo.c.FindOne(ctx, bson.M{"_id": objectId}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, storage_service.ErrNewsNotFound
		}

		return nil, err
	}

	return doc.ToNews(), nil
}

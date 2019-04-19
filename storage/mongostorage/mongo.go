package mongostorage

import (
	"context"
	"fmt"
	"time"

	"github.com/matrosov-nikita/newsapp/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/matrosov-nikita/newsapp/storage/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NewsMongoDocument struct {
	ID        primitive.ObjectID `bson:"_id"`
	Header    string             `bson:"header"`
	CreatedAt time.Time          `bson:"created_at"`
}

func NewsDocument(news *pb.News) *NewsMongoDocument {
	t, _ := ptypes.Timestamp(news.CreatedAt)
	return &NewsMongoDocument{ID: primitive.NewObjectID(), Header: news.Header, CreatedAt: t}
}

func (d *NewsMongoDocument) ToNews() *pb.News {
	createdAt, _ := ptypes.TimestampProto(d.CreatedAt)
	return &pb.News{
		Id:        d.ID.Hex(),
		Header:    d.Header,
		CreatedAt: createdAt,
	}
}

type NewsRepository struct {
	c *mongo.Collection
}

func CreateNewsRepository(client *mongo.Client) *NewsRepository {
	return &NewsRepository{
		c: client.Database("newsdb").Collection("news"),
	}
}

func (repo *NewsRepository) Save(news *pb.News) (string, error) {
	doc := NewsDocument(news)
	_, err := repo.c.InsertOne(context.Background(), doc)
	if err != nil {
		return "", err
	}
	return doc.ID.Hex(), nil
}

func (repo *NewsRepository) FindByID(id string) (*pb.News, error) {
	ctx := context.Background()
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("incorrect id: %s", id)
	}

	var doc NewsMongoDocument
	if err := repo.c.FindOne(ctx, bson.M{"_id": objectId}).Decode(&doc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, storage.ErrNewsNotFound
		}

		return nil, err
	}

	return doc.ToNews(), nil
}

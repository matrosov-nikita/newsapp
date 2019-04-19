package main

import (
	"context"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/matrosov-nikita/newsapp/storage"
	"github.com/matrosov-nikita/newsapp/storage/mongostorage"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/nats-io/go-nats"
)

func main() {
	natsURL := getEnv("NATS_URL", nats.DefaultURL)
	natsConnection, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("could not connect to NATS broker: %v", err)
	}

	mongoDSN := getEnv("MONGO_DSN", "mongodb://localhost:27017")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDSN))
	if err != nil {
		log.Fatalf("could not connect to mongodb server: %v", err)
	}

	newsStorage := mongostorage.CreateNewsRepository(client)
	st := storage.NewStorageService(newsStorage)
	subs := NewSubs(st, natsConnection)
	_, _ = natsConnection.Subscribe("news.create", subs.CreateNews)
	_, _ = natsConnection.Subscribe("news.get", subs.FindNews)

	runtime.Goexit()
}

func getEnv(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	return defaultVal
}

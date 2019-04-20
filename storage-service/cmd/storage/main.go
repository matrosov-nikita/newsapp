package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/matrosov-nikita/newsapp/storage-service"
	"github.com/matrosov-nikita/newsapp/storage-service/mongo-storage"
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

	newsStorage := mongo_storage.CreateNewsRepository(client)
	st := storage_service.NewStorageService(newsStorage)
	subs := NewSubs(st, natsConnection)
	natsConnection.Subscribe("news.create", subs.CreateNews)
	natsConnection.Subscribe("news.get", subs.FindNews)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)

	<-shutdownCh
	log.Println("Gracefully stopping...")
	natsConnection.Close()
	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("fail when close mongodb client:%v\n", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	os.Exit(0)
}

func getEnv(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	return defaultVal
}

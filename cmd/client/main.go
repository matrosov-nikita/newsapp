package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/matrosov-nikita/newsapp/newsclient/nats"

	"github.com/matrosov-nikita/newsapp/newsclient"

	"github.com/gorilla/mux"
)

func main() {
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")
	natsClient, err := nats.New(natsURL)
	if err != nil {
		log.Fatalf("could not connect to NATS broker: %v", err)
	}

	client := newsclient.New(natsClient)
	r := mux.NewRouter()
	h := NewHandler(client)
	h.Attach(r)

	addr := getEnv("SERVER_ADDRESS", ":8888")
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Println("Listening on", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)

	<-shutdownCh
	log.Println("Gracefully stopping...")
	natsClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}

func getEnv(name string, defaultVal string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}

	return defaultVal
}

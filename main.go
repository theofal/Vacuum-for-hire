package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func getDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	client, err := pubsub.NewClient(context.Background(), getDotEnvVar("PROJECT_ID"), option.WithAPIKey(getDotEnvVar("API_KEY")))
	if err != nil {
		fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()
	// Use the authenticated client.

}

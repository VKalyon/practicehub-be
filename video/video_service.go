package video

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

//encore:service
type VideoService struct {
	mongoUri string
	client   *mongo.Client
}

func initVideoService() (*VideoService, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	service := VideoService{
		mongoUri: mongodbURI,
	}

	service.connect()

	return &service, nil
}

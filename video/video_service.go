package video

import (
	"go.mongodb.org/mongo-driver/mongo"
)

//encore:service
type VideoService struct {
	client *mongo.Client
}

var secrets struct {
	MONGODB_URI string
}

func initVideoService() (*VideoService, error) {
	go measureMemory()

	service := VideoService{}

	service.connect()

	return &service, nil
}

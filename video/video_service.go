package video

import "go.mongodb.org/mongo-driver/mongo"

//encore:service
type VideoService struct {
	mongoUri string
	client   *mongo.Client
}

func initVideoService() (*VideoService, error) {
	service := VideoService{
		mongoUri: "mongodb://localhost:27017",
	}

	service.connect()

	return &service, nil
}

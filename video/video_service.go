package video

import (
	//_ "net/http/pprof"

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
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	go measureMemory()

	service := VideoService{}

	service.connect()

	return &service, nil
}

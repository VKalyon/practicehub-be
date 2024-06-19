package video

import (
	//_ "net/http/pprof"

	"encore.app/middleware"
	"encore.dev/pubsub"
	"go.mongodb.org/mongo-driver/mongo"
)

//encore:service
type VideoService struct {
	client *mongo.Client
}

var secrets struct {
	MONGODB_URI string
}

var videoUploadedTopicRef = pubsub.TopicRef[pubsub.Publisher[*middleware.VideoUploadedEvent]](middleware.VideoUploaded)

func initVideoService() (*VideoService, error) {
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	go measureMemory()

	service := VideoService{}

	service.connect()

	return &service, nil
}

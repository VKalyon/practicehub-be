package video

import (
	"bytes"
	"context"
	"os"

	"encore.dev/rlog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *VideoService) connect() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(secrets.MONGODB_URI).SetServerAPIOptions(serverAPI)

	var err error
	s.client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}
}

func (s *VideoService) uploadVideo(video *os.File) (primitive.ObjectID, error) {
	uploadOpts := options.GridFSUpload()

	db := s.client.Database("practiceHubVideo")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return primitive.NilObjectID, err
	}

	objectID, err := bucket.UploadFromStream(video.Name(), video, uploadOpts)
	if err != nil {
		return objectID, err
	}

	rlog.Info("New file uploaded with ID %s", objectID)

	return objectID, err
}

func (s *VideoService) getVideoById(mongoidHex string) {
	db := s.client.Database("practiceHubVideo")
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		rlog.Error(err.Error())
		return
	}

	id, err := primitive.ObjectIDFromHex(mongoidHex)
	if err != nil {
		rlog.Error(err.Error())

		return
	}

	fileBuffer := bytes.NewBuffer(nil)
	if _, err := bucket.DownloadToStream(id, fileBuffer); err != nil {
		rlog.Error(err.Error())
	}

	rlog.Info(string(fileBuffer.Len()))
}

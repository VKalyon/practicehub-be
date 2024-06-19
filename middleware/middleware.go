package middleware

import "encore.dev/pubsub"

type VideoUploadedEvent struct {
	Message string
}

var VideoUploaded = pubsub.NewTopic[*VideoUploadedEvent]("video-uploaded", pubsub.TopicConfig{DeliveryGuarantee: pubsub.AtLeastOnce})

package notification

import (
	"context"
	"fmt"

	"encore.app/middleware"
	"encore.dev/pubsub"
	"encore.dev/rlog"
)

//encore:service
type NotificationService struct {
}

var Subscription = pubsub.NewSubscription(middleware.VideoUploaded, "notification-service", pubsub.SubscriptionConfig[*middleware.VideoUploadedEvent]{
	Handler:        HandleEvent,
	RetryPolicy:    &pubsub.RetryPolicy{MaxRetries: 10},
	MaxConcurrency: 20,
})

func HandleEvent(ctx context.Context, msg *middleware.VideoUploadedEvent) error {
	s := fmt.Sprintf("Received message: %v from topic: %s", msg, "notification-service")
	rlog.Info(s)
	return nil
}

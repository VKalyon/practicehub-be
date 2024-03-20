package video

import (
	"context"
)

//encore:api public path=/video/:id
func GetVideo(ctx context.Context, id int) (*Response, error) {
	response := Response{URL: "test", Id: id, Title: "testTitle"}

	return &response, nil
}

type Response struct {
	URL   string
	Id    int
	Title string
}

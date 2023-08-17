package httpadapter

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ConstructFetchContent() func(ctx context.Context, url string) (io.ReadCloser, error) {
	return func(ctx context.Context, url string) (io.ReadCloser, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			log.Println("ERROR", "creating request", err.Error())
			return nil, fmt.Errorf("something went wrong")
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("ERROR", "request failing", err.Error())
			return nil, fmt.Errorf("error when fetching content")
		}

		return resp.Body, nil
	}
}

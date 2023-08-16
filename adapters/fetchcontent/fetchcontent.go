package fetchcontent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Construct() func(ctx context.Context, url string) (*bytes.Reader, error) {
	return func(ctx context.Context, url string) (*bytes.Reader, error) {
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

		bodyBts, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("ERROR", "cannot read response", err.Error())
			return nil, fmt.Errorf("error when fetching content")
		}
		resp.Body.Close()

		return bytes.NewReader(bodyBts), nil
	}
}

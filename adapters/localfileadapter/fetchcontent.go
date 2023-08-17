package localfileadapter

import (
	"context"
	"fmt"
	"io"
	"os"
)

func ConstructFetchContent() func(ctx context.Context, url string) (io.ReadCloser, error) {
	return (&fetchcontent{}).Execute
}

type fetchcontent struct {
}

func (e *fetchcontent) Execute(ctx context.Context, url string) (io.ReadCloser, error) {
	file, err := os.Open(url)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not existed")
		}
		return nil, err
	}
	return file, nil
}

package adapters

import (
	"context"
	"io"
)

type FetchContent func(ctx context.Context, url string) (io.ReadCloser, error)
type StoreContent func(ctx context.Context, path string, io io.Reader) error

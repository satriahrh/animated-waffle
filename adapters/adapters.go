package adapters

import (
	"bytes"
	"context"
)

type FetchContent func(ctx context.Context, url string) (*bytes.Reader, error)
type StoreContent func(ctx context.Context, path string, reader *bytes.Reader) error

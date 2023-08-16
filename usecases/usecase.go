package usecases

import "context"

type FetchUrl func(ctx context.Context, url string) error

package usecases

import (
	"context"

	"github.com/satriahrh/animated-waffle/entity"
)

type FetchUrl func(ctx context.Context, url string) error
type FetchMetadata func(ctx context.Context, url string) (entity.Metadata, error)

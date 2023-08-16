package command

import (
	"context"
	"time"

	"github.com/satriahrh/autify-tht/adapters"
	"github.com/satriahrh/autify-tht/adapters/fetchcontent"
	"github.com/satriahrh/autify-tht/adapters/storecontent"
	"github.com/satriahrh/autify-tht/sanitizer"
	"github.com/satriahrh/autify-tht/usecases"
	"github.com/satriahrh/autify-tht/usecases/fetchurl"
	"github.com/spf13/cobra"
)

var Command *cobra.Command

func init() {
	Command = &cobra.Command{
		Use:  "fetch {set of urls}",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), sanitizer.SanitizeFetch),
		RunE: func(cmd *cobra.Command, urls []string) error {
			ctx, cancelFunc := context.WithTimeout(cmd.Context(), time.Minute)
			defer cancelFunc()

			var fetchContent adapters.FetchContent = fetchcontent.Construct()
			var storeContent adapters.StoreContent = storecontent.Construct()
			var usecase usecases.FetchUrl = fetchurl.Construct(fetchContent, storeContent)

			for _, url := range urls {
				ctx, cancelFunc := context.WithCancel(ctx)
				defer cancelFunc()

				err := usecase(ctx, url)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

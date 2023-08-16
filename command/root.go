package command

import (
	"context"
	"fmt"
	"sync"
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

			var wg sync.WaitGroup

			duplicatePreventionLookup := make(map[string]struct{})
			for _, url := range urls {
				if _, isDuplicated := duplicatePreventionLookup[url]; isDuplicated {
					continue
				}
				duplicatePreventionLookup[url] = struct{}{}

				wg.Add(1)
				go func(url string) {
					var usecase usecases.FetchUrl = fetchurl.Construct(fetchContent, storeContent)
					ctx, cancelFunc := context.WithCancel(ctx)
					defer cancelFunc()

					err := usecase(ctx, url)
					if err != nil {
						fmt.Println(time.Now().Local().String(), url, "ERROR", err.Error())
					} else {
						fmt.Println(time.Now().Local().String(), url, "SUCCESS")
					}
					wg.Done()
				}(url)

			}
			wg.Wait()
			return nil
		},
	}
}

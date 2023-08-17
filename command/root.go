package command

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/satriahrh/animated-waffle/adapters/httpadapter"
	"github.com/satriahrh/animated-waffle/adapters/localfileadapter"
	"github.com/satriahrh/animated-waffle/sanitizer"
	"github.com/satriahrh/animated-waffle/usecases"
	"github.com/satriahrh/animated-waffle/usecases/fetchmetadata"
	"github.com/satriahrh/animated-waffle/usecases/fetchurl"
	"github.com/spf13/cobra"
)

var Command *cobra.Command
var metadataFlagged bool

func init() {
	Command = &cobra.Command{
		Use:  "fetch {set of urls}",
		Args: cobra.MatchAll(cobra.MinimumNArgs(1), sanitizer.SanitizeFetch),
		RunE: fetchHandler,
	}
	Command.PersistentFlags().BoolVarP(&metadataFlagged, "metadata", "", false, "Print the metadata of what was fetched.")
}

func fetchHandler(cmd *cobra.Command, urls []string) error {
	ctx, cancelFunc := context.WithTimeout(cmd.Context(), time.Minute)
	defer cancelFunc()

	var wg sync.WaitGroup

	duplicatePreventionLookup := make(map[string]struct{})
	for _, url := range urls {
		if _, isDuplicated := duplicatePreventionLookup[url]; isDuplicated {
			continue
		}
		duplicatePreventionLookup[url] = struct{}{}

		wg.Add(1)
		go fetchHandlerSingleUrl(ctx, &wg, url)
	}
	wg.Wait()
	return nil
}

func fetchHandlerSingleUrl(ctx context.Context, wg *sync.WaitGroup, url string) error {
	defer wg.Done()

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	var fetchurl usecases.FetchUrl = fetchurl.Construct(
		httpadapter.ConstructFetchContent(),
		localfileadapter.ConstructStoreContent(),
	)

printingMetadata:
	if metadataFlagged {
		var fetchmetadata usecases.FetchMetadata = fetchmetadata.Construct(
			localfileadapter.ConstructFetchContent(),
		)
		metadata, err := fetchmetadata(ctx, url)
		if err != nil {
			if err.Error() != "no metadata found" {
				return fmt.Errorf("ERROR METADATA %s: %s", url, err.Error())
			}
		} else {
			metadata.FmtPrintln()
			return nil
		}
	}

	err := fetchurl(ctx, url)
	if err != nil {
		return fmt.Errorf("ERROR FETCHING %s: %s", url, err.Error())
	}

	if metadataFlagged {
		goto printingMetadata
	}
	return nil
}

package fetchmetadata

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/satriahrh/autify-tht/adapters"
	"github.com/satriahrh/autify-tht/entity"
	"github.com/satriahrh/autify-tht/utils"
)

func Construct(fetchContent adapters.FetchContent) func(ctx context.Context, url string) (entity.Metadata, error) {
	return (&executor{fetchContent}).Execute
}

type executor struct {
	fetchContent adapters.FetchContent
}

func (e *executor) Execute(ctx context.Context, url string) (metadata entity.Metadata, err error) {
	metadataFilepath := utils.UrlToMetadataFilepath(url)

	reader, err := e.fetchContent(ctx, metadataFilepath)
	if err != nil {
		if err.Error() == "not existed" {
			return entity.Metadata{}, fmt.Errorf("no metadata found")
		}
		log.Println("ERROR", "cannot read metadata html content", err.Error())
		return entity.Metadata{}, fmt.Errorf("problem with reading metadata")
	}
	defer reader.Close()

	err = json.NewDecoder(reader).Decode(&metadata)
	if err != nil {
		log.Println("ERROR", "cannot parse metadata file", err.Error())
		return entity.Metadata{}, fmt.Errorf("problem with reading metadata")
	}

	return
}

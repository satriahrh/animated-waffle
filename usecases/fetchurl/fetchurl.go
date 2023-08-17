package fetchurl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/satriahrh/autify-tht/adapters"
	"github.com/satriahrh/autify-tht/entity"
	"github.com/satriahrh/autify-tht/utils"
	"golang.org/x/net/html"
)

func Construct(fetchContent adapters.FetchContent, storeContent adapters.StoreContent) func(ctx context.Context, url string) error {
	return (&executor{fetchContent, storeContent}).Execute
}

type executor struct {
	fetchContent adapters.FetchContent
	storeContent adapters.StoreContent
}

func (e *executor) Execute(ctx context.Context, url string) error {
	reader, err := e.fetchContent(ctx, url)
	if err != nil {
		return err
	}
	defer reader.Close()

	filepath := utils.UrlToFilepath(url)

	htmlBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Println("ERROR", "cannot read all html contnt", err.Error())
		return fmt.Errorf("problem with reading html content")
	}

	err = e.storeContent(ctx, filepath+".html", bytes.NewReader(htmlBytes))
	if err != nil {
		log.Println("ERROR", "cannot store html content", err.Error())
		return fmt.Errorf("problem with storing html")
	}

	metadata := generateMetadata(ctx, url, bytes.NewBuffer(htmlBytes))

	metadataFilepath := utils.UrlToMetadataFilepath(url)
	err = e.storeContent(ctx, metadataFilepath, metadata.JsonEncodeToIoReader())
	if err != nil {
		log.Println("ERROR", "cannot store metadata html content", err.Error())
		return fmt.Errorf("problem with storing metadata html")
	}

	return nil
}

func generateMetadata(ctx context.Context, urlStr string, reader io.Reader) (metadata entity.Metadata) {
	mp := struct {
		*entity.Metadata
		sync.Mutex
	}{
		&metadata, sync.Mutex{},
	}

	mp.URL = urlStr
	mp.LastFetch = time.Now()

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Println("ERROR", "url cannot be parsed", err.Error())
		return
	}
	mp.Site = parsedURL.Host

	htmlDoc, err := html.Parse(reader)
	if err != nil {
		log.Println("ERROR", "content cannot be parsed with html", err.Error())
		return
	}

	var wg sync.WaitGroup
	var processingNode = func(node *html.Node) {

		if node.Type == html.ElementNode {
			switch node.Data {
			case "a":
				mp.Lock()
				mp.NumLinks += 1
				mp.Unlock()
			case "img":
				mp.Lock()
				mp.Images += 1
				mp.Unlock()
			}
		}

		wg.Done()
	}

	queue := []*html.Node{htmlDoc}
	wg.Add(1)
	for len(queue) != 0 {
		node := queue[0]
		go processingNode(node)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			wg.Add(1)
			queue = append(queue, child)
		}
		queue = queue[1:]
	}

	wg.Wait()
	return
}

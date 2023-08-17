package fetchurl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"path"
	"strings"
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

	metadata, updatedHtmlReader, err := e.processingHtml(ctx, url, bytes.NewBuffer(htmlBytes))
	if err != nil {
		log.Println("ERROR", "problem with processing html: ", err.Error())
		return nil // no error for now
		// return fmt.Errorf("problem with processing html")
	}

	err = e.storeContent(ctx, filepath+".html", updatedHtmlReader)
	if err != nil {
		log.Println("ERROR", "cannot store html content", err.Error())
		return fmt.Errorf("problem with storing html")
	}

	metadataFilepath := utils.UrlToMetadataFilepath(url)
	err = e.storeContent(ctx, metadataFilepath, metadata.JsonEncodeToIoReader())
	if err != nil {
		log.Println("ERROR", "cannot store metadata html content", err.Error())
		return fmt.Errorf("problem with storing metadata html")
	}

	return nil
}

type htmlProcessorData struct {
	*sync.Mutex
	*sync.WaitGroup
	*entity.Metadata
	*url.URL
}

func (e *executor) processingHtml(ctx context.Context, urlStr string, reader io.Reader) (entity.Metadata, io.Reader, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		log.Println("ERROR", "url cannot be parsed: ", err.Error())
		return entity.Metadata{}, nil, fmt.Errorf("something is not supposed to be happened")
	}

	mp := &htmlProcessorData{
		&sync.Mutex{}, &sync.WaitGroup{}, &entity.Metadata{}, parsedURL,
	}

	mp.Metadata.URL = urlStr
	mp.Metadata.LastFetch = time.Now()

	mp.Metadata.Site = parsedURL.Host

	htmlDoc, err := html.Parse(reader)
	if err != nil {
		log.Println("ERROR", "content cannot be parsed with html: ", err.Error())
		return *mp.Metadata, nil, fmt.Errorf("not a valid html file")
	}

	queue := []*html.Node{htmlDoc}
	mp.WaitGroup.Add(1)
	for len(queue) != 0 {
		node := queue[0]
		go e.processingHtmlNode(ctx, mp, node)
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			mp.WaitGroup.Add(1)
			queue = append(queue, child)
		}
		queue = queue[1:]
	}

	mp.WaitGroup.Wait()

	var buffer strings.Builder
	err = html.Render(&buffer, htmlDoc)
	if err != nil {
		return *mp.Metadata, nil, fmt.Errorf("asset downloaded but html is not updated")
	}

	return *mp.Metadata, strings.NewReader(buffer.String()), nil
}

func (e *executor) processingHtmlNode(ctx context.Context, mp *htmlProcessorData, node *html.Node) {
	defer mp.WaitGroup.Done()

	var sourceToRetrieve string
	var keyToModify string

	if node.Type == html.ElementNode {
		switch node.Data {
		case "a":
			mp.Mutex.Lock()
			mp.Metadata.NumLinks += 1
			mp.Mutex.Unlock()
		case "img":
			mp.Mutex.Lock()
			mp.Metadata.Images += 1
			mp.Mutex.Unlock()
		}
		if node.Data == "img" || node.Data == "script" || node.Data == "source" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					sourceToRetrieve = attr.Val
					keyToModify = "src"
				}
			}
		} else if node.Data == "link" {
			var isCss bool
			var source string
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					source = attr.Val
				} else if attr.Key == "rel" && (attr.Val == "stylesheet" || attr.Val == "icon") {
					isCss = true
				}
			}
			if isCss {
				sourceToRetrieve = source
				keyToModify = "href"
			}
		}
	}

	if sourceToRetrieve != "" {
		assetUrl := sourceToRetrieve
		if !utils.UrlRegex.MatchString(sourceToRetrieve) {
			assetUrl = fmt.Sprintf("%s://%s%s", mp.URL.Scheme, mp.URL.Host, sourceToRetrieve)
		}
		reader, err := e.fetchContent(ctx, assetUrl)
		if err != nil {
			log.Println("ERROR", "cannot fetch asset content: ", assetUrl, ": ", err.Error())
			return
		}
		defer reader.Close()

		assetPath := utils.UrlToFilepath(assetUrl)
		assetCompletePath := path.Join("."+utils.UrlToFilepath(mp.Metadata.URL), assetPath)
		err = e.storeContent(ctx, assetCompletePath, reader)
		if err != nil {
			log.Println("ERROR", "cannot store asset content: ", assetUrl, ": ", err.Error())
			return
		}

		for i, attr := range node.Attr {
			if attr.Key == keyToModify {
				node.Attr[i].Val = assetCompletePath
				break
			}
		}
	}
}

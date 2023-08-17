package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Metadata struct {
	URL       string    `json:"url"`
	Site      string    `json:"site"`
	NumLinks  int       `json:"num_links"`
	Images    int       `json:"images"`
	LastFetch time.Time `json:"last_fetch"`
}

func (v Metadata) JsonEncodeToIoReader() io.Reader {
	jsonBytes := bytes.NewBuffer(nil)
	json.NewEncoder(jsonBytes).Encode(v)
	return jsonBytes
}

func (v Metadata) FmtPrintln() {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("*** " + v.URL + " ***\n")
	buf.WriteString("site: " + v.Site + "\n")
	buf.WriteString("num_links: " + fmt.Sprint(v.NumLinks) + "\n")
	buf.WriteString("images: " + fmt.Sprint(v.Images) + "\n")
	buf.WriteString("last_fetch: " + v.LastFetch.Local().String() + "\n")
	fmt.Println(buf.String())
}

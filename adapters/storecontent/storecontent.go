package storecontent

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
)

func Construct() func(ctx context.Context, path string, reader *bytes.Reader) error {
	return func(ctx context.Context, path string, reader *bytes.Reader) error {
		file, err := os.Create(path)
		if err != nil {
			log.Println("ERROR", "cannot create file", err.Error())
			return fmt.Errorf("cannot creating file")
		}
		defer file.Close()

		_, err = reader.WriteTo(file)
		if err != nil {
			log.Println("ERROR", "cannot write file", err.Error())
			return fmt.Errorf("cannot writing to file")
		}

		return nil
	}
}

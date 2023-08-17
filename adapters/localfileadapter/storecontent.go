package localfileadapter

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func ConstructStoreContent() func(ctx context.Context, filepath string, reader io.Reader) error {
	return func(ctx context.Context, filepath string, reader io.Reader) error {
		os.Mkdir(path.Dir(filepath), 0755)
		file, err := os.Create(filepath)
		if err != nil {
			log.Println("ERROR", "cannot create file", err.Error())
			return fmt.Errorf("cannot creating file")
		}
		defer file.Close()

		bt := make([]byte, 100)
		n, err := reader.Read(bt)
		for ; n > 0 && err == nil; n, err = reader.Read(bt) {
			_, err := file.Write(bt[:n])
			if err != nil {
				log.Println("ERROR", "cannot write file", err.Error())
				return fmt.Errorf("cannot writing to file")
			}
		}
		if err != nil && err != io.EOF {
			log.Println("ERROR", "cannot read content", err.Error())
			return fmt.Errorf("cannot read content")
		}
		return nil
	}
}

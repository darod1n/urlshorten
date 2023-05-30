package storage

import (
	"flag"
	"log"
	"os"
)

var path string

func init() {

	flag.StringVar(&path, "f", "/tmp/short-url-db.json", "File path")

	if fp := os.Getenv("FILE_STORAGE_PATH"); fp != "" {
		path = fp
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
	}

}

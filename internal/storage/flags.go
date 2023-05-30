package storage

import (
	"flag"
	"os"
)

var path string

func init() {

	flag.StringVar(&path, "f", "tmp/short-url-db.json", "File path")

	if fp := os.Getenv("FILE_STORAGE_PATH"); fp != "" {
		path = fp
	}
}

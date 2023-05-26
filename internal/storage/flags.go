package storage

import "flag"

var path string

func init() {
	flag.StringVar(&path, "f", "tmp/test.json", "File path")
}

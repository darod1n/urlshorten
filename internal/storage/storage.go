package storage

type DB struct {
	urls map[string]string
}

type Storage interface {
	AddURL(url string, shortURL string)
	GetURL(shortURL string) (string, bool)
}

func (db *DB) AddURL(url string, shortURL string) {
	db.urls[shortURL] = url
}

func (db *DB) GetURL(shortURL string) (string, bool) {
	bigURL, ok := db.urls[shortURL]
	return bigURL, ok
}

func NewDB() *DB {
	return &DB{urls: map[string]string{}}
}

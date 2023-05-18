package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

type Storage interface{}

func WithLoggin(s Storage, w http.ResponseWriter, res *http.Request) {
	start := time.Now()

	uri := res.RequestURI
	method := res.Method

	duration := time.Since(start)

	sugar.Infow(
		"uri", uri,
		"method", method,
		"duration", duration,
	)
}

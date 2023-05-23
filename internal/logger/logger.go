package logger

import (
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

var Sugar zap.SugaredLogger

func InitializeLoger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	Sugar = *logger.Sugar()
	return nil
}

func WithLoggin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, res *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		uri := res.RequestURI
		method := res.Method

		h.ServeHTTP(&lw, res)

		duration := time.Since(start)
		body, _ := io.ReadAll(res.Body)

		Sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)

		Sugar.Infoln(
			"uri", uri,
			"method", method,
			"status", responseData.status,
			"body", body,
		)
	})

}

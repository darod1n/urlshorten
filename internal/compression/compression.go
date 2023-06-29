package compression

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{zw: gzip.NewWriter(w), ResponseWriter: w}
}

const (
	contentEncodingKey  = "Content-Encoding"
	gzipKey             = "gzip"
	typeApplicationJSON = "application/json"
	typeTextHTML        = "text/html"
	statusCodesGZIP     = 300
)

func (gz *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < statusCodesGZIP {
		gz.ResponseWriter.Header().Set(contentEncodingKey, gzipKey)
	}
	gz.ResponseWriter.WriteHeader(statusCode)
}

func (gz *gzipWriter) Write(b []byte) (int, error) {
	return gz.zw.Write(b)
}

func (gz *gzipWriter) Close() error {
	return gz.zw.Close()
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{r: r, zr: zr}, nil
}

func (gz *gzipReader) Close() error {
	if err := gz.r.Close(); err != nil {
		return err
	}
	return gz.zr.Close()
}

func (gz gzipReader) Read(p []byte) (n int, err error) {
	return gz.zr.Read(p)
}

func WithCompress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentType := r.Header.Get("Content-Type")
		supportGzip := strings.Contains(acceptEncoding, gzipKey) && (strings.Contains(contentType, typeApplicationJSON) || strings.Contains(contentType, typeTextHTML))
		if supportGzip {
			gw := newGzipWriter(w)
			ow = gw

			defer gw.Close()
		}

		contentEncoding := r.Header.Get(contentEncodingKey)

		sendsGzip := strings.Contains(contentEncoding, gzipKey)
		if sendsGzip {
			gr, err := newGzipReader(r.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = gr

			defer gr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

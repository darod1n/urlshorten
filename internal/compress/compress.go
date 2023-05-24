package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{w: w, zw: gzip.NewWriter(w)}
}

func (gz *gzipWriter) Header() http.Header {
	return gz.w.Header()
}

func (gz *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		gz.w.Header().Set("Content-Encoding", "gzip")

	}
	gz.w.WriteHeader(statusCode)
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

func (gz gzipReader) Close() error {
	if err := gz.r.Close(); err != nil {
		return err
	}
	return gz.zr.Close()
}

func (gz *gzipReader) Read(b []byte) (n int, err error) {
	return gz.zr.Read(b)
}

func WithCompress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		if supportGzip {
			gw := newGzipWriter(w)

			ow = gw

			defer gw.Close()

		}

		contentEncoding := r.Header.Get("Content-Encoding")
		contentType := r.Header.Get("Content-Type")

		sendsGzip := strings.Contains(contentEncoding, "gzip") && (strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/html"))
		if sendsGzip {
			gr, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = gr

			defer gr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}

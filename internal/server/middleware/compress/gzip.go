package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (w compressWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (c *compressReader) Read(p []byte) (int, error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func Compressgzip() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				h.ServeHTTP(w, r)
				return
			}

			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")

			h.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
		})
	}
}

func Decompressgzip() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
				h.ServeHTTP(w, r)
				return
			}

			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer cr.Close()
			r.Body = cr
			h.ServeHTTP(w, r)
		})
	}
}

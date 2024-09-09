package server

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
		srv.logger.Errorf("Ошибка считывания данных: %v", err)
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func Compressgzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			srv.logger.Errorf("Недопустимый уровень сжатия: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func Decompressgzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.logger.Debugf("Header: %v", r.Header)
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		cr, err := newCompressReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer cr.Close()
		r.Body = cr
		next.ServeHTTP(w, r)
	})
}

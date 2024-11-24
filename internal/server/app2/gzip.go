package server

import (
	"compress/gzip"
	"fmt"
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
	srv.logger.Debugf("Сжатие данных...")
	srv.logger.Debugf("Сжатие данных прошло успешно")
	return w.Writer.Write(b)
}

func (c *compressReader) Read(p []byte) (int, error) {
	srv.logger.Debugf("Чтение данных...")
	srv.logger.Debugf("Чтение данных прошло успешно")
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
	fmt.Println("Работа Compressgzip...")
	defer fmt.Println("Работа Compressgzip завершена")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Проверка, что клиент поддерживает формат сжатых данных gzip...")
		fmt.Println("Получение заголовка Accept-Encoding...")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fmt.Println("Заголовок не содержит Accept-Encoding: gzip")
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("Заголовок содержит Accept-Encoding: gzip")

		fmt.Println("Создание gzip.Writer c уровнем сжатия BestSpeed...")
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			//		srv.logger.Errorf("Недопустимый уровень сжатия: %v", err)
			fmt.Printf("Недопустимый уровень сжатия: %v", err)
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("Создание gzip.Writer c уровнем сжатия BestSpeed прошло успешно")

		defer gz.Close()

		fmt.Println("Установка заголовка Content-Encoding: gzip...")
		w.Header().Set("Content-Encoding", "gzip")
		fmt.Println("Установка заголовка Content-Encoding: gzip прошла успешно")

		fmt.Println("Передача хендлеру переменной типа gzipWriter для вывода данных")
		next.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func Decompressgzip(next http.Handler) http.Handler {
	fmt.Println("Работа Decompressgzip...")
	defer fmt.Println("Работа Decompressgzip завершена")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Получение заголовка Content-Encoding...")
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			fmt.Println("Заголовок не содержит Content-Encoding: gzip")
			next.ServeHTTP(w, r)
			return
		}

		fmt.Println("Заголовок содержит Content-Encoding: gzip")

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

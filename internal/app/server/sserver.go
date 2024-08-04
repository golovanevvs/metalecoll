package server

import (
	"io"
	"net/http"
	"strings"
)

type server struct {
	router *http.ServeMux
}

type Metric struct {
	Name    string
	gauge   float64
	counter int64
}

type MemStorage struct {
	metrics map[string]Metric
}

type MetStorage interface {
	Add(metric Metric) error
}

func newServer() *server {
	s := &server{
		router: http.NewServeMux(),
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("", s.handleMainFunc())
}

func (s *server) handleMainFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sbody := strings.Split(string(body), "/")
		if len(sbody) != 5 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mT := sbody[2]
		mN := sbody[3]
		mV := sbody[4]
	}
}

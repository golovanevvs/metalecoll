package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

type server struct {
	store  mapstorage.Storage
	router *chi.Mux
}

var srv *server

func Start() {
	store := mapstorage.NewStorage()
	srv = NewServer(store)
	fmt.Println("Запущен сервер:", constants.Addr)
	err := http.ListenAndServe(constants.Addr, srv)
	if err != nil {
		fmt.Println("Ошибка сервера")
		panic(err)
	}
}

func NewServer(store mapstorage.Storage) *server {
	s := &server{
		store:  store,
		router: chi.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Post("/update/{metType}/{metName}/{matValue}", MainHandle)
	s.router.Get("/", GetMetricNamesHandle)
	s.router.Get("/value/{metType}/{metName}", GetMetricValueHandle)
}

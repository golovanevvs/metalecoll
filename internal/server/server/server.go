package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
)

type server struct {
	store  storage.Storage
	router *chi.Mux
}

var srv *server

func Start(config *Config) {
	store := mapstorage.NewStorage()
	srv = NewServer(store, config)
	fmt.Println("Запущен сервер:", config.Addr)
	err := http.ListenAndServe(config.Addr, srv)
	if err != nil {
		fmt.Println("Ошибка сервера")
		panic(err)
	}
}

func NewServer(store storage.Storage, config *Config) *server {
	s := &server{
		store:  store,
		router: chi.NewRouter(),
	}
	s.configureRouter(config)
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter(config *Config) {
	str := fmt.Sprintf("/%s/{%s}/{%s}/{%s}", config.UpdateMethod, constants.MetTypeURL, constants.MetNameURL, constants.MetValueURL)
	s.router.Post(str, func(w http.ResponseWriter, r *http.Request) {
		UpdateMetricsHandler(w, r, s.store)
	})
	s.router.Get("/", GetMetricNamesHandler)
	str = fmt.Sprintf("/%s/{%s}/{%s}", config.GetValueMethod, constants.MetTypeURL, constants.MetNameURL)
	s.router.Get(str, GetMetricValueHandler)
}

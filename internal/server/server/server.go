package server

import (
	"flag"
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
var flagRunAddr string

func Start() {
	parseFlags()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	store := mapstorage.NewStorage()
	srv = NewServer(store)
	fmt.Println("Запущен сервер:", flagRunAddr)
	err := http.ListenAndServe(flagRunAddr, srv)
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

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", constants.Addr, "address and port to run server")
	flag.Parse()
}

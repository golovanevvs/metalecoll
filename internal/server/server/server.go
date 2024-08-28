package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"go.uber.org/zap"
)

type server struct {
	store  storage.Storage
	router *chi.Mux
	logger *zap.SugaredLogger
}

var srv *server

func Start(config *Config) {
	store := mapstorage.NewStorage()
	srv = NewServer(store, config)
	//fmt.Println("Запущен сервер:", config.Addr)
	srv.logger.Infof("Запущен сервер: %s", config.Addr)
	if err := http.ListenAndServe(config.Addr, srv); err != nil {
		srv.logger.Fatalf(err.Error(), "event", "start server")
	}
}

func NewServer(store storage.Storage, config *Config) *server {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize zap")
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	s := &server{
		store:  store,
		router: chi.NewRouter(),
		logger: sugar,
	}

	s.configureRouter(config)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

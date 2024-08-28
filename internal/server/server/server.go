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
	logger *zap.Logger
}

var srv *server

func Start(config *Config) {
	store := mapstorage.NewStorage()
	srv = NewServer(store, config)
	//fmt.Println("Запущен сервер:", config.Addr)
	srv.logger.Info("Запущен сервер: ", zap.String("Addr", config.Addr))
	if err := http.ListenAndServe(config.Addr, srv); err != nil {
		srv.logger.Fatal("Ошибка запуска сервера", zap.Error(err))
	}
}

func NewServer(store storage.Storage, config *Config) *server {
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	panic("cannot initialize zap")
	// }
	// defer logger.Sync()

	//sugar := logger.Sugar()

	log, err := Initialize("info")
	if err != nil {
		panic("cannot initialize zap")
	}

	s := &server{
		store:  store,
		router: chi.NewRouter(),
		logger: log,
	}

	s.configureRouter(config)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

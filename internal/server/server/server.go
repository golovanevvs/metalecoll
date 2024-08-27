package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
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

func (s *server) configureRouter(config *Config) {
	str := fmt.Sprintf("/%s/{%s}/{%s}/{%s}",
		config.UpdateMethod,
		constants.MetTypeURL,
		constants.MetNameURL,
		constants.MetValueURL,
	)
	s.router.Post(str, func(w http.ResponseWriter, r *http.Request) {
		UpdateMetricsHandler(w, r, s.store)
	})

	s.router.Get("/", GetMetricNamesHandler)
	str = fmt.Sprintf("/%s/{%s}/{%s}",
		config.GetValueMethod,
		constants.MetTypeURL,
		constants.MetNameURL,
	)
	s.router.Get(str, GetMetricValueHandler)
}

// func WithLogging(h http.Handler) http.Handler {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {
// 		h.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(logFn)
// }

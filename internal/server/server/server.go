package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"github.com/sirupsen/logrus"
)

type server struct {
	store  storage.Storage
	router *chi.Mux
	//logger *zap.Logger
	logger *logrus.Logger
}

var srv *server

func Start(config *Config) {
	store := mapstorage.NewStorage()

	srv = NewServer(store, config)

	//fmt.Println("Запущен сервер:", config.Addr)

	//srv.logger.Info("Запущен сервер: ", zap.String("Addr", config.Addr))

	srv.logger.Infof("Запущен сервер: %s", config.Addr)

	if err := http.ListenAndServe(config.Addr, srv); err != nil {
		//srv.logger.Fatal("Ошибка запуска сервера", zap.Error(err))
		srv.logger.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func NewServer(store storage.Storage, config *Config) *server {
	// logger, err := zap.NewDevelopment()
	// if err != nil {
	// 	panic("cannot initialize zap")
	// }
	// defer logger.Sync()

	//sugar := logger.Sugar()

	// log, err := InitializeLogger("info")
	// if err != nil {
	// 	panic("cannot initialize zap")
	// }

	logLogrus := logrus.New()
	logLogrus.SetLevel(logrus.DebugLevel)
	logLogrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	s := &server{
		store:  store,
		router: chi.NewRouter(),
		logger: logLogrus,
	}

	s.configureRouter(config)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

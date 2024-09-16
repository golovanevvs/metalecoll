package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/filestorage"
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

	//db, err := newDB(config.DatabaseDNS)
	// if err != nil {
	// 	return err
	// }

	//defer db.Close()

	srv = NewServer(store, config)

	if config.Restore {
		srv.logger.Debugf("Восстановление метрик из файла %v...", config.FileStoragePath)
		err := filestorage.GetFromFile(config.FileStoragePath, srv.store)
		if err != nil {
			srv.logger.Errorf("Ошибка чтения данных из файла: %v. Сервер будет запущен без восстановления метрик", err)
		} else {
			srv.logger.Infof("Восстановление метрик из файла %v прошло успешно", config.FileStoragePath)
		}
	}

	//srv.logger.Info("Запущен сервер: ", zap.String("Addr", config.Addr))

	go func() {
		srv.logger.Infof("Запущен сервер: %s", config.Addr)
		if err := http.ListenAndServe(config.Addr, srv); err != nil {
			//srv.logger.Fatal("Ошибка запуска сервера", zap.Error(err))
			srv.logger.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	saveIntTime := time.NewTicker(time.Duration(config.StoreInterval) * time.Second)
	defer saveIntTime.Stop()

	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-saveIntTime.C:
				srv.logger.Debugf("Сохранение метрик в файл %v...", config.FileStoragePath)
				if err := filestorage.SaveToFile(config.FileStoragePath, srv.store); err != nil {
					srv.logger.Errorf("Ошибка сохранения в файл: %v", err)
				}
			case <-stop:
				srv.logger.Debugf("Стоп")
			}
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	srv.logger.Infof("Завершение работы сервера...")

	srv.logger.Debugf("Сохранение метрик в файл %v...", config.FileStoragePath)
	if err := filestorage.SaveToFile(config.FileStoragePath, srv.store); err != nil {
		srv.logger.Errorf("Ошибка сохранения в файл: %v", err)
	}

	srv.logger.Infof("Работа сервера завершена. Всем спасибо.")
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
	l, _ := logrus.ParseLevel((config.LogLevel))
	logLogrus.SetLevel(l)
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

// func newDB(databaseDNS string) (*sql.DB, error) {
// 	db, err := sql.Open("pgx", databaseDNS)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

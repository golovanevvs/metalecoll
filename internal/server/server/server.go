package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/filestorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/sqlstorage"
	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type server struct {
	store   storage.Storage
	storeDB storage.StorageDB
	router  *chi.Mux
	//logger *zap.Logger
	logger *logrus.Logger
}

var srv *server

func Start(config *Config) error {
	store := mapstorage.NewStorage()

	db, err := newDB(config.DatabaseDNS)
	if err != nil {
		return err
	}

	defer db.Close()

	storeDB := sqlstorage.New(db)

	srv = NewServer(store, storeDB, config)

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
	return nil
}

func NewServer(store storage.Storage, storeDB storage.StorageDB, config *Config) *server {
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
		store:   store,
		storeDB: storeDB,
		router:  chi.NewRouter(),
		logger:  logLogrus,
	}

	s.configureRouter(config)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newDB(databaseDNS string) (*sql.DB, error) {
	tableMetrics := "metrics"

	db, err := sql.Open("pgx", databaseDNS)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	exist, err := tablesExist(ctx, db, tableMetrics)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return nil, err
	}

	if !exist {
		fmt.Printf("Создание таблицы %v...\n", tableMetrics)
		_, err = db.ExecContext(ctx, "CREATE TABLE metrics (id INTEGER PRIMARY KEY)")
		if err != nil {
			fmt.Printf("Ошибка создания таблицы %v: %v\n", tableMetrics, err)
			return nil, err
		}
		exist2, err := tablesExist(ctx, db, tableMetrics)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
			return nil, err
		}
		if !exist2 {
			return nil, errors.New("неизвестная ошибка создания табюлицы metrics")
		}
		fmt.Printf("Создание таблицы %v прошло успешно\n", tableMetrics)
	}

	return db, nil
}

func tablesExist(ctx context.Context, db *sql.DB, nameTable string) (bool, error) {
	var exists bool

	fmt.Printf("Проверка, что таблица %v существует...\n", nameTable)

	row := db.QueryRowContext(ctx,
		"SELECT EXISTS "+
			"(SELECT FROM information_schema.tables "+
			"WHERE table_name = 'metrics')")

	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Printf("Таблица %v существует\n", nameTable)
		return true, nil
	}
	fmt.Printf("Таблицы %v не существует\n", nameTable)
	return false, nil
}

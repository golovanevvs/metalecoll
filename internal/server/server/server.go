package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/mapstorage"
	_ "github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type server struct {
	mapStore mapstorage.Storage
	dbStore  storage.StorageDB
	router   *chi.Mux
	//logger *zap.Logger
	logger *logrus.Logger
	c      *config.Config
}

var srv *server

// Start запускает сервер
func Start(c *config.Config) error {
	// Инициализация map-хранилища
	mapStore := mapstorage.New()

	// Выбор и инициализация основного хранилища: если флаг (-d) пуст, то выбирается файловое хранилище, иначе - БД
	dbStore, err := storage.New(c)
	if err != nil {
		fmt.Printf("Ошибка инициализации хранилища: %v\n", err)
		return err
	}

	// Инициализация сервера
	srv = NewServer(mapStore, dbStore, c)

	// Вывод информации об основном хранилище
	nameDB := srv.dbStore.GetNameDB()
	srv.logger.Infof("Основное хранилище: %v", nameDB)

	// Инициализация контекста Background
	ctx := context.Background()

	// Восстановление метрик в map-хранилище из БД при запуске сервера, если флаг (-r) = true
	if c.Storage.Restore {
		srv.logger.Debugf("Восстановление метрик из основного хранилища %v...", nameDB)
		interFromDB, err := srv.dbStore.GetMetricsFromDB(ctx, c)
		if err != nil {
			srv.logger.Errorf("Ошибка чтения данных из основного хранилища: %v. Сервер будет запущен без восстановления метрик", err)
		} else {
			mapStore.Metrics = interFromDB.GetMetrics()
			srv.logger.Infof("Восстановление метрик из основного хранилища %v прошло успешно", nameDB)
		}
	}

	//srv.logger.Info("Запущен сервер: ", zap.String("Addr", config.Addr))

	// Запуск сервера
	go func() {
		srv.logger.Infof("Запущен сервер: %s", c.Server.Addr)
		if err := http.ListenAndServe(c.Server.Addr, srv); err != nil {
			//srv.logger.Fatal("Ошибка запуска сервера", zap.Error(err))
			srv.logger.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Сохранение метрик из map-хранилища в основное хранилище через интервал StoreInterval (-i)
	saveIntTime := time.NewTicker(time.Duration(c.Server.StoreInterval) * time.Second)
	defer saveIntTime.Stop()

	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-saveIntTime.C:
				srv.logger.Debugf("Сохранение метрик в основное хранилище %v...", nameDB)
				if err := srv.dbStore.SaveMetricsToDB(ctx, c, srv.mapStore); err != nil {
					srv.logger.Errorf("Ошибка сохранения метрик в основное хранилище: %v", err)
				} else {
					srv.logger.Debugf("Сохранение метрик в основное хранилище прошло успешно")
				}
			case <-stop:
				srv.logger.Debugf("Стоп")
			}
		}
	}()

	// Сохранение метрик из map-хранилища в основное хранилище при завершении работы сервера
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	srv.logger.Infof("Завершение работы сервера...")

	srv.logger.Debugf("Сохранение метрик в основное хранилище %v...", nameDB)
	if err := srv.dbStore.SaveMetricsToDB(ctx, c, srv.mapStore); err != nil {
		srv.logger.Errorf("Ошибка сохранения метрик в основное хранилище: %v", err)
	}

	srv.logger.Infof("Работа сервера завершена. Всем спасибо.")
	return nil
}

// NewServer - конструктор сервера
func NewServer(mapStore mapstorage.Storage, dbStore storage.StorageDB, c *config.Config) *server {
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
	l, _ := logrus.ParseLevel(c.Logger.LogLevel)
	logLogrus.SetLevel(l)
	logLogrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	s := &server{
		mapStore: mapStore,
		dbStore:  dbStore,
		router:   chi.NewRouter(),
		logger:   logLogrus,
		c:        c,
	}

	s.configureRouter(c)

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// func newDB(databaseDNS string) (*sql.DB, error) {
// 	tableMetrics := "metrics"

// 	db, err := sql.Open("pgx", databaseDNS)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	exist, err := tablesExist(ctx, db, tableMetrics)
// 	if err != nil {
// 		fmt.Printf("Ошибка: %v\n", err)
// 		return nil, err
// 	}

// 	if !exist {
// 		fmt.Printf("Создание таблицы %v...\n", tableMetrics)
// 		_, err = db.ExecContext(ctx, "CREATE TABLE metrics (id INTEGER PRIMARY KEY)")
// 		if err != nil {
// 			fmt.Printf("Ошибка создания таблицы %v: %v\n", tableMetrics, err)
// 			return nil, err
// 		}
// 		exist2, err := tablesExist(ctx, db, tableMetrics)
// 		if err != nil {
// 			fmt.Printf("Ошибка: %v\n", err)
// 			return nil, err
// 		}
// 		if !exist2 {
// 			return nil, errors.New("неизвестная ошибка создания табюлицы metrics")
// 		}
// 		fmt.Printf("Создание таблицы %v прошло успешно\n", tableMetrics)
// 	}

// 	return db, nil
// }

// func tablesExist(ctx context.Context, db *sql.DB, nameTable string) (bool, error) {
// 	var exists bool

// 	fmt.Printf("Проверка, что таблица %v существует...\n", nameTable)

// 	row := db.QueryRowContext(ctx,
// 		"SELECT EXISTS "+
// 			"(SELECT FROM information_schema.tables "+
// 			"WHERE table_name = 'metrics')")

// 	err := row.Scan(&exists)
// 	if err != nil {
// 		return false, err
// 	}

// 	if exists {
// 		fmt.Printf("Таблица %v существует\n", nameTable)
// 		return true, nil
// 	}
// 	fmt.Printf("Таблицы %v не существует\n", nameTable)
// 	return false, nil
// }

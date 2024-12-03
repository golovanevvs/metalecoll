// Модуль app предназначен для запуска сервера и корректного завершения его работы.
package app

import (
	"context"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/filestorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/postgres"
	"github.com/golovanevvs/metalecoll/internal/server/transport/http/handler"
	"github.com/sirupsen/logrus"
)

// RunApp запускает сервер и корректно завершает его работу.
func RunApp() {
	//! подготовительные операции
	// инициализация логгера
	lg := logrus.New()

	// инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		lg.Fatalf("ошибка инициализации конфигурации сервера: %s", err.Error())
	}

	// установка уровня логгирования
	lg.SetLevel(cfg.Logger.LogLevel)

	// Выбор и инициализация основного хранилища: если флаг (-d) пуст, то выбирается файловое хранилище, иначе - БД
	var mainStore storage.IStorageDB
	switch cfg.Storage.DatabaseDSN {
	case "":
		mainStore = filestorage.NewFileStorage(cfg.Storage.FileStoragePath)
	default:
		mainStore, err = postgres.NewPostgres(cfg.Storage.DatabaseDSN)
		if err != nil {
			lg.Fatalf("Ошибка инициализации базы данных: %s", err.Error())
		}
	}

	// инициализация map-хранилища
	mst := mapstorage.NewMapStorage()
	// инициализация хранилища
	st := storage.NewStorage(mainStore)
	// инициализация сервиса
	sv := service.NewService(mst, st)
	// инициализация хендлера
	hd := handler.NewHandler(sv, lg, cfg.Crypto.HashKey)
	// инициализация сервера
	srv := NewServer()

	// Инициализация контекста Background
	ctx := context.Background()

	// Вывод информации об основном хранилище
	nameDB := st.GetNameDB()
	lg.Infof("Основное хранилище: %v", nameDB)

	// Восстановление метрик в map-хранилище из БД при запуске сервера, если флаг (-r) = true
	if cfg.Storage.Restore {
		lg.Debugf("Восстановление метрик из основного хранилища %s...", nameDB)
		interFromDB, err := st.GetMetricsFromDB(ctx, cfg)
		if err != nil {
			lg.Errorf("Ошибка чтения данных из основного хранилища: %s. Сервер будет запущен без восстановления метрик", err.Error())
		} else {
			mst.Metrics = interFromDB.GetMetricsFromMap()
			lg.Infof("Восстановление метрик из основного хранилища %v прошло успешно", nameDB)
		}
	}

	//! запуск сервера
	go func() {
		lg.Infof("Сервер сбора метрик metalecoll запущен")
		if err := srv.RunServer(cfg.Server.Addr, hd.InitRoutes()); err != nil {
			lg.Fatalf("работа сервера: %s", err.Error())
		}
	}()

	//! запуск профилировщика
	// go func() {
	// 	lg.Infof("Сервер профилировщика запущен")
	// 	if err := http.ListenAndServe(":9090", nil); err != nil {
	// 		lg.Fatalf("ошибка запуска сервера профилировщика: %s", err.Error())
	// 	}
	// }()

	// Сохранение метрик из map-хранилища в основное хранилище через интервал StoreInterval (-i)
	saveIntTime := time.NewTicker(time.Duration(cfg.Server.StoreInterval) * time.Second)
	defer saveIntTime.Stop()

	stop := make(chan bool)

	//! сохранение данных в основное хранилище
	go func() {
		for {
			select {
			case <-saveIntTime.C:
				lg.Debugf("Сохранение метрик в основное хранилище %s...", nameDB)
				if err := st.SaveMetricsToDB(ctx, cfg, mst); err != nil {
					lg.Errorf("Ошибка сохранения метрик в основное хранилище: %s", err.Error())
				} else {
					lg.Debugf("Сохранение метрик в основное хранилище прошло успешно")
				}
			case <-stop:
				lg.Debugf("Сохранение метрик в основное хранилище остановлено")
				//TODO требуется доработка обработки сигнала
				return
			}
		}
	}()

	//! операции при завершении работы сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	lg.Infof("Получен сигнал о завершении работы сервера")

	// Сохранение метрик из map-хранилища в основное хранилище при завершении работы сервера
	lg.Debugf("Сохранение метрик в основное хранилище %s...", nameDB)
	if err := st.SaveMetricsToDB(ctx, cfg, mst); err != nil {
		lg.Errorf("Ошибка сохранения метрик в основное хранилище: %v", err)
	}

	if err := srv.ShutdownServer(context.Background()); err != nil {
		lg.Errorf("Ошибка при завершении работы сервера: %s", err.Error())
	}

	if err := st.CloseDB(); err != nil {
		lg.Errorf("Ошибка при завершении работы с БД: %s", err.Error())
	}

	lg.Infof("Работа сервера корректно завершена")
	time.Sleep(time.Second * 2)
}

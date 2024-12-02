package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/golovanevvs/metalecoll/internal/server/storage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/filestorage"
	"github.com/golovanevvs/metalecoll/internal/server/storage/postgres"
	"github.com/sirupsen/logrus"
)

func TestUpdateMetric(t *testing.T) {
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
	hd := NewHandler(sv, lg, cfg.Crypto.HashKey)
	// инициализация сервера

	//! Задание входных и выходных параметров
	type actual struct {
		targetRequest string
	}
	type want struct {
		httpstatus int
	}
	tests := []struct {
		name   string
		actual actual
		want   want
	}{
		{
			name: "positive test #1",
			actual: actual{
				targetRequest: "/update/gauge/gauge1/0.1",
			},
			want: want{
				httpstatus: http.StatusOK,
			},
		},
	}

	//! Запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// создание имитации запроса
			request := httptest.NewRequest(http.MethodPost, test.actual.targetRequest, nil)

			// создание нового Recorder
			w := httptest.NewRecorder()

			// запуск хендлера
			hd.UpdateMetric(w, request)

			//

			res := w.Result()
		})
	}
}

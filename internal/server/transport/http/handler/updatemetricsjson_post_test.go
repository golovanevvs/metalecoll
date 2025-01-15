package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/dto"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/mocks"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricsJSON(t *testing.T) {
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

	// инициализация map-хранилища
	mst := mapstorage.NewMapStorage()

	//! использование заглушки БД
	// создание контроллера gomock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создание объекта-заглушки
	m := mocks.NewMockIStorageDB(ctrl)

	// инициализация сервиса
	sv := service.NewService(mst, m)
	// инициализация хендлера
	hd := NewHandler(sv, lg, cfg.Crypto.HashKey, cfg.Crypto.PrivateKeyPath)

	// инициализация тестового сервера
	ts := httptest.NewServer(hd.InitRoutes())
	defer ts.Close()

	//! Задание входных и выходных параметров
	// входные значения
	type actual struct {
		body dto.Metrics // данные для тела запроса
	}

	// ожидаемые значения
	type want struct {
		httpStatus int         // код ответа
		resp       dto.Metrics // ответ сервера
	}

	// параметры тестов
	tests := []struct {
		name   string // имя теста
		actual actual // входные значения
		want   want   // ожидаемые значения
	}{
		{
			name: "positive test gauge: gauge1=0.1",
			actual: actual{
				body: dto.Metrics{
					ID:    "gauge1",
					MType: "gauge",
					Delta: nil,
					Value: floatToPtr(0.1),
				},
			},
			want: want{
				httpStatus: http.StatusOK,
				resp: dto.Metrics{
					ID:    "gauge1",
					MType: "gauge",
					Delta: nil,
					Value: floatToPtr(0.1),
				},
			},
		},
		{
			name: "positive test update gauge value: gauge1=0.2",
			actual: actual{
				body: dto.Metrics{
					ID:    "gauge1",
					MType: "gauge",
					Delta: nil,
					Value: floatToPtr(0.2),
				},
			},
			want: want{
				httpStatus: http.StatusOK,
				resp: dto.Metrics{
					ID:    "gauge1",
					MType: "gauge",
					Delta: nil,
					Value: floatToPtr(0.2),
				},
			},
		},
		{
			name: "positive test: counter1=1",
			actual: actual{
				body: dto.Metrics{
					ID:    "counter1",
					MType: "counter",
					Delta: intToPtr(1),
					Value: nil,
				},
			},
			want: want{
				httpStatus: http.StatusOK,
				resp: dto.Metrics{
					ID:    "counter1",
					MType: "counter",
					Delta: intToPtr(1),
					Value: nil,
				},
			},
		},
		{
			name: "positive test update counter value: counter1=2",
			actual: actual{
				body: dto.Metrics{
					ID:    "counter1",
					MType: "counter",
					Delta: intToPtr(2),
					Value: nil,
				},
			},
			want: want{
				httpStatus: http.StatusOK,
				resp: dto.Metrics{
					ID:    "counter1",
					MType: "counter",
					Delta: intToPtr(3),
					Value: nil,
				},
			},
		},
		{
			name: "negative test: no name",
			actual: actual{
				body: dto.Metrics{
					ID:    "",
					MType: "gauge",
					Delta: nil,
					Value: floatToPtr(0.1),
				},
			},
			want: want{
				httpStatus: http.StatusNotFound,
				resp: dto.Metrics{
					ID:    "",
					MType: "gauge",
					Delta: nil,
					Value: floatToPtr(0.1),
				},
			},
		},
		{
			name: "negative test: wrong type",
			actual: actual{
				body: dto.Metrics{
					ID:    "counter3",
					MType: "wrongtype",
					Delta: intToPtr(4),
					Value: nil,
				},
			},
			want: want{
				httpStatus: http.StatusBadRequest,
				resp: dto.Metrics{
					ID:    "counter3",
					MType: "wrongtype",
					Delta: intToPtr(4),
					Value: nil,
				},
			},
		},
	}

	//! Запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// сериализация тела запроса
			encAct, err := json.Marshal(test.actual.body)
			if err != nil {
				t.Fatalf("ошибка сериализации тела запроса: %s", err.Error())
			}
			reqBody := bytes.NewReader(encAct)

			// отправка запроса и формирование ответа и тела ответа
			resp, respBody := testRequest(t, ts, "POST", "/update/", reqBody)
			defer resp.Body.Close()

			// сериализация ожидаемого значения для сравнения с ответом
			encExp, err := json.Marshal(test.want.resp)
			if err != nil {
				t.Fatalf("ошибка сериализации ожидаемого тела ответа: %s", err.Error())
			}

			// сравнение ответа и ожидаемого значения
			assert.Equal(t, test.want.httpStatus, resp.StatusCode)
			if test.want.httpStatus == http.StatusOK {
				assert.Equal(t, string(encExp), respBody)
			}
		})
	}
}

func intToPtr(i int64) *int64 {
	return &i
}

func floatToPtr(f float64) *float64 {
	return &f
}

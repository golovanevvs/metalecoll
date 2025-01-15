package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
	"github.com/golovanevvs/metalecoll/internal/server/mocks"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
		targetRequest string // запрос
	}

	// ожидаемые значения
	type want struct {
		httpStatus int    // код ответа
		resp       string // ответ сервера
	}

	// параметры тестов
	tests := []struct {
		name   string // имя теста
		actual actual // входные значения
		want   want   // ожидаемые значения
	}{
		{
			name: "positive test gauge: /update/gauge/gauge1/0.1",
			actual: actual{
				targetRequest: "/update/gauge/gauge1/0.1",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: gauge, name: gauge1, value: 0.1",
			},
		},
		{
			name: "positive test update gauge: /update/gauge/gauge1/0.2",
			actual: actual{
				targetRequest: "/update/gauge/gauge1/0.2",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: gauge, name: gauge1, value: 0.2",
			},
		},
		{
			name: "positive test counter: /update/counter/counter1/1",
			actual: actual{
				targetRequest: "/update/counter/counter1/1",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: counter, name: counter1, value: 1",
			},
		},
		{
			name: "positive test update counter: /update/counter/counter1/2",
			actual: actual{
				targetRequest: "/update/counter/counter1/2",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: counter, name: counter1, value: 3",
			},
		},
		{
			name: "negative test: no name /update/gauge//0.1",
			actual: actual{
				targetRequest: "/update/gauge//0.1",
			},
			want: want{
				httpStatus: http.StatusNotFound,
				resp:       "",
			},
		},
		{
			name: "negative test: wrong gauge value /update/gauge/gauge3/x",
			actual: actual{
				targetRequest: "/update/gauge/gauge3/x",
			},
			want: want{
				httpStatus: http.StatusBadRequest,
				resp:       "",
			},
		},
		{
			name: "negative test: wrong counter value /update/counter/counter3/0.1",
			actual: actual{
				targetRequest: "/update/counter/counter3/0.1",
			},
			want: want{
				httpStatus: http.StatusBadRequest,
				resp:       "",
			},
		},
		{
			name: "negative test: wrong type /update/wrongtype/counter3/4",
			actual: actual{
				targetRequest: "/update/wrongtype/counter3/4",
			},
			want: want{
				httpStatus: http.StatusBadRequest,
				resp:       "",
			},
		},
	}

	//! Запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, respBody := testRequest(t, ts, "POST", test.actual.targetRequest, nil)
			defer resp.Body.Close()
			assert.Equal(t, test.want.httpStatus, resp.StatusCode)
			if test.want.httpStatus == http.StatusOK {
				assert.Equal(t, test.want.resp, respBody)
			}
		})
	}
}

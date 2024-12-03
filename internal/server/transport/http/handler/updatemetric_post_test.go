package handler

import (
	"io"
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
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestHandler(t *testing.T) {
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
	// создание контроллера
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создание объекта-заглушки
	m := mocks.NewMockIStorageDB(ctrl)

	// инициализация сервиса
	sv := service.NewService(mst, m)
	// инициализация хендлера
	hd := NewHandler(sv, lg, cfg.Crypto.HashKey)

	ts := httptest.NewServer(hd.InitRoutes())
	defer ts.Close()

	//! Задание входных и выходных параметров
	type actual struct {
		targetRequest string
	}
	type want struct {
		httpStatus int
		resp       string
	}
	tests := []struct {
		name   string
		actual actual
		want   want
	}{
		{
			name: "positive test /update/gauge/gauge1/0.1",
			actual: actual{
				targetRequest: "/update/gauge/gauge1/0.1",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: gauge, name: gauge1, value: 0.1",
			},
		},
		{
			name: "positive test /update/gauge/gauge2/0.2",
			actual: actual{
				targetRequest: "/update/gauge/gauge1/0.2",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: gauge, name: gauge1, value: 0.2",
			},
		},
		{
			name: "positive test /update/counter/counter1/1",
			actual: actual{
				targetRequest: "/update/counter/counter1/1",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: counter, name: counter1, value: 1",
			},
		},
		{
			name: "positive test /update/counter/counter1/2",
			actual: actual{
				targetRequest: "/update/counter/counter1/2",
			},
			want: want{
				httpStatus: http.StatusOK,
				resp:       "type: counter, name: counter1, value: 3",
			},
		},
	}

	//! Запуск тестов
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, respBody := testRequest(t, ts, "POST", test.actual.targetRequest)
			assert.Equal(t, test.want.httpStatus, resp.StatusCode)
			assert.Equal(t, test.want.resp, respBody)
		})
	}
}

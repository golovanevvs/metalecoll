package handler

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/golovanevvs/metalecoll/internal/server/config"
// 	"github.com/golovanevvs/metalecoll/internal/server/dto"
// 	"github.com/golovanevvs/metalecoll/internal/server/mapstorage"
// 	"github.com/golovanevvs/metalecoll/internal/server/mocks"
// 	"github.com/golovanevvs/metalecoll/internal/server/service"
// 	"github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/assert"
// )

// func TestUpdateMetricJSON(t *testing.T) {
// 	//! подготовительные операции
// 	// инициализация логгера
// 	lg := logrus.New()

// 	// инициализация конфигурации
// 	cfg, err := config.NewConfig()
// 	if err != nil {
// 		lg.Fatalf("ошибка инициализации конфигурации сервера: %s", err.Error())
// 	}

// 	// установка уровня логгирования
// 	lg.SetLevel(cfg.Logger.LogLevel)

// 	// инициализация map-хранилища
// 	mst := mapstorage.NewMapStorage()

// 	//! использование заглушки БД
// 	// создание контроллера gomock
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// создание объекта-заглушки
// 	m := mocks.NewMockIStorageDB(ctrl)

// 	// инициализация сервиса
// 	sv := service.NewService(mst, m)
// 	// инициализация хендлера
// 	hd := NewHandler(sv, lg, cfg.Crypto.HashKey)

// 	// инициализация тестового сервера
// 	ts := httptest.NewServer(hd.InitRoutes())
// 	defer ts.Close()

// 	//! Задание входных и выходных параметров
// 	// входные значения
// 	type actual struct {
// 		body string // тело запроса
// 	}

// 	// ожидаемые значения
// 	type want struct {
// 		httpStatus int         // код ответа
// 		resp       dto.Metrics // ответ сервера
// 	}

// 	// параметры тестов
// 	tests := []struct {
// 		name   string // имя теста
// 		actual actual // входные значения
// 		want   want   // ожидаемые значения
// 	}{
// 		{
// 			name: "positive test: /update/gauge/gauge1/0.1",
// 			actual: actual{
// 				body: `{
// 					"id": "gauge1",
// 					"type": "gauge",
// 					"delta":,
// 					"value": 0.1
// 				}`,
// 			},
// 			want: want{
// 				httpStatus: http.StatusOK,
// 				resp: dto.Metrics{
// 					ID: "gauge1",
// 					MType: "gauge",
// 					Delta: i,
// 					Value: 0.1,
// 				},
// 			},
// 		},
// 	}

// 	//! Запуск тестов
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			// сериализация тела запроса
// 			encAct, err := json.Marshal(test.actual.body)
// 			if err != nil {
// 				t.Fatalf("ошибка сериализации тела запроса: %s", err.Error())
// 			}
// 			reqBody := bytes.NewReader(encAct)

// 			// формирование тела ответа
// 			resp, respBody := testRequest(t, ts, "POST", "/update/", reqBody)
// 			defer resp.Body.Close()

// 			a, _ := io.ReadAll(reqBody)
// 			fmt.Printf("reqBody: %s\n", string(a))
// 			fmt.Printf("respBody: %s\n", respBody)
// 			// var respJSON dto.Metrics
// 			// decAct := json.NewDecoder(resp.Body)
// 			// if err := decAct.Decode(&respJSON); err != nil {
// 			// 	t.Fatalf("ошибка десериализации тела ответа: %s", err.Error())
// 			// }

// 			encExp, err := json.Marshal(test.want.resp)
// 			if err != nil {
// 				t.Fatalf("ошибка сериализации ожидаемого тела ответа: %s", err.Error())
// 			}

// 			fmt.Println(string(encExp))
// 			assert.Equal(t, test.want.httpStatus, resp.StatusCode)
// 			assert.Equal(t, string(encExp), respBody)
// 		})
// 	}
// }

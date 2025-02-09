// Модуль handler предназначен для обработки запросов.
package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/middleware/compress"
	"github.com/golovanevvs/metalecoll/internal/server/middleware/decrypt"
	"github.com/golovanevvs/metalecoll/internal/server/middleware/logger"
	"github.com/golovanevvs/metalecoll/internal/server/middleware/trustedipchecker"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// структура handler
type handler struct {
	sv             *service.Service
	lg             *logrus.Logger
	hashKey        string
	privateKeyPath string
	trustedSubnet  string
}

// NewHandler - конструктор *handler.
func NewHandler(sv *service.Service, lg *logrus.Logger, hashKey string, privateKeyPath string, trustedSubnet string) *handler {
	return &handler{
		sv:             sv,
		lg:             lg,
		hashKey:        hashKey,
		privateKeyPath: privateKeyPath,
		trustedSubnet:  trustedSubnet,
	}
}

// InitRoutes - маршрутизация запросов, используется в качестве http.Handler при запуске сервера.
func (hd *handler) InitRoutes() *chi.Mux {
	// создание экземпляра роутера
	rt := chi.NewRouter()

	// использование middleware
	// логгирование
	rt.Use(logger.WithLogging(hd.lg))
	// компрессия
	rt.Use(compress.Compressgzip())
	// декомпрессия
	rt.Use(compress.Decompressgzip())

	// маршруты
	rt.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", hd.UpdateMetric)
		r.With(
			func(next http.Handler) http.Handler {
				return decrypt.Decrypt(hd.privateKeyPath, hd.lg, next)
			},
			func(next http.Handler) http.Handler {
				return trustedipchecker.TrustedIPChecker(hd.trustedSubnet, next)
			},
		).Post("/", hd.UpdateMetricsJSON)
	})
	rt.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", hd.GetMetricValue)
		r.Post("/", hd.GetMetricValueJSON)
	})
	rt.Post("/updates/", hd.UpdatesMetricsJSON)
	rt.Get("/ping", hd.PingDatabase)
	rt.Get("/", hd.GetMetricNames)

	// маршруты для профилировщика
	rtProf := chi.NewRouter()

	rtProf.HandleFunc("/", pprof.Index)
	rtProf.HandleFunc("/cmdline", pprof.Cmdline)
	rtProf.HandleFunc("/profile", pprof.Profile)
	rtProf.HandleFunc("/symbol", pprof.Symbol)
	rtProf.HandleFunc("/trace", pprof.Trace)

	rtProf.Handle("/goroutine", pprof.Handler("goroutine"))
	rtProf.Handle("/heap", pprof.Handler("heap"))
	rtProf.Handle("/threadcreate", pprof.Handler("threadcreate"))
	rtProf.Handle("/block", pprof.Handler("block"))
	rtProf.Handle("/mutex", pprof.Handler("mutex"))
	rtProf.Handle("/allocs", pprof.Handler("allocs"))

	rt.Mount("/debug/pprof", rtProf)

	//rtProf.HandleFunc("/debug/pprof/", pprof.Index)
	return rt
}

// testRequest формирует запрос, отправляет запрос на тестовый сервер, фозвращает ответ от сервера и тело ответа.
func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	// формирование запроса
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	// отправка запроса на тестовый сервер и получение ответа
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// формирование тела ответа тестового сервера
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

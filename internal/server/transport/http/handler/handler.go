package handler

import (
	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/middleware/compress"
	"github.com/golovanevvs/metalecoll/internal/server/middleware/logger"
	"github.com/golovanevvs/metalecoll/internal/server/service"
	"github.com/sirupsen/logrus"
)

// структура handler
type handler struct {
	sv      *service.Service
	lg      *logrus.Logger
	hashKey string
}

// NewHandler - конструктор *handler
func NewHandler(sv *service.Service, lg *logrus.Logger, hashKey string) *handler {
	return &handler{
		sv:      sv,
		lg:      lg,
		hashKey: hashKey,
	}
}

// InitRoutes - маршрутизация запросов, используется в качестве http.Handler при запуске сервера
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
		r.Post("/", hd.UpdateMetricsJSON)
	})
	rt.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", hd.GetMetricValue)
		r.Post("/", hd.GetMetricValueJSON)
	})
	rt.Post("/updates/", hd.UpdatesMetricsJSON)
	rt.Get("/ping", hd.PingDatabase)
	rt.Get("/", hd.GetMetricNames)

	return rt
}

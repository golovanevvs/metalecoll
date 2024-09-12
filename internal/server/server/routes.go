package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *server) configureRouter(config *Config) {
	s.router.Use(func(h http.Handler) http.Handler {
		return WithLogging(h)
	})

	s.router.Use(func(h http.Handler) http.Handler {
		return Compressgzip(h)
	})

	s.router.Use(func(h http.Handler) http.Handler {
		return Decompressgzip(h)
	})

	s.router.Route(fmt.Sprintf("/%s", config.UpdateMethod), func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск UpdateMetricsHandler")
			UpdateMetricsHandler(w, r, s.store)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск UpdateMetricsJSONHandler")
			UpdateMetricsJSONHandler(w, r, s.store)
		})
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		srv.logger.Debugf("Запуск GetMetricNamesHandler")
		GetMetricNamesHandler(w, r, s.store)
	})

	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		srv.logger.Debugf("Запуск PingDatabaseHandler")
		PingDatabaseHandler(w, r, config.DatabaseDNS)
	})

	s.router.Route(fmt.Sprintf("/%s", config.GetValueMethod), func(r chi.Router) {
		r.Get("/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск GetMetricValueHandler")
			GetMetricValueHandler(w, r, s.store)

		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск GetMetricValueJSONHandler")
			GetMetricValueJSONHandler(w, r, s.store)

		})

	})
}

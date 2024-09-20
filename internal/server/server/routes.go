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
			s.logger.Debugf("Запуск UpdateMetricsHandler")
			s.UpdateMetricsHandler(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			s.logger.Debugf("Запуск UpdateMetricsJSONHandler")
			s.UpdateMetricsJSONHandler(w, r)
		})
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("Запуск GetMetricNamesHandler")
		s.GetMetricNamesHandler(w, r)
	})

	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("Запуск PingDatabaseHandler")
		s.PingDatabaseHandler(w, r)
	})

	s.router.Route(fmt.Sprintf("/%s", config.GetValueMethod), func(r chi.Router) {
		r.Get("/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
			s.logger.Debugf("Запуск GetMetricValueHandler")
			s.GetMetricValueHandler(w, r)

		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			s.logger.Debugf("Запуск GetMetricValueJSONHandler")
			s.GetMetricValueJSONHandler(w, r)

		})

	})
}

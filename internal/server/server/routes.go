package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/config"
)

func (s *server) configureRouter(c *config.Config) {
	s.router.Use(func(h http.Handler) http.Handler {
		return WithLogging(h)
	})

	s.router.Use(func(h http.Handler) http.Handler {
		return Compressgzip(h)
	})

	s.router.Use(func(h http.Handler) http.Handler {
		return Decompressgzip(h)
	})

	s.router.Route("/update", func(r chi.Router) {
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

	s.router.Post("/updates", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debugf("Запуск UpdatesMetricsJSONHandler")
		s.UpdatesMetricsJSONHandler(w, r)
	})

	s.router.Route("/value", func(r chi.Router) {
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

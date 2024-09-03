package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

func (s *server) configureRouter(config *Config) {
	var str string

	s.router.Use(func(h http.Handler) http.Handler {
		return WithLogging(h)
	})

	s.router.Use(func(h http.Handler) http.Handler {
		return Compressgzip(h)
	})

	s.router.Use(func(h http.Handler) http.Handler {
		return Decompressgzip(h)
	})

	str = fmt.Sprintf("/{%s}/{%s}/{%s}",
		constants.MetTypeURL,
		constants.MetNameURL,
		constants.MetValueURL,
	)

	s.router.Route(fmt.Sprintf("/%s", config.UpdateMethod), func(r chi.Router) {
		r.Post(str, func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск UpdateMetricsHandler")
			UpdateMetricsHandler(w, r, s.store)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск UpdateMetricsJSONHandler")
			UpdateMetricsJSONHandler(w, r, s.store)
		})
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		GetMetricNamesHandler(w, r, s.store)
	})

	str = fmt.Sprintf("/{%s}/{%s}",
		constants.MetTypeURL,
		constants.MetNameURL,
	)
	s.router.Route(fmt.Sprintf("/%s", config.GetValueMethod), func(r chi.Router) {
		r.Get(str, func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск GetMetricValueHandler")
			GetMetricValueHandler(w, r, s.store)

		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			srv.logger.Debugf("Запуск GetMetricValueJSONHandler")
			GetMetricValueJSONHandler(w, r, s.store)

		})

	})
}

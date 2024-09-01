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

	str = fmt.Sprintf("/{%s}/{%s}/{%s}",
		constants.MetTypeURL,
		constants.MetNameURL,
		constants.MetValueURL,
	)

	s.router.Route(fmt.Sprintf("/%s", constants.UpdateMethod), func(r chi.Router) {
		r.Post(str, func(w http.ResponseWriter, r *http.Request) {
			UpdateMetricsHandler(w, r, s.store)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			UpdateMetricsJSONHandler(w, r, s.store)
		})
	})

	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		GetMetricNamesHandler(w, r, s.store)
	})

	str = fmt.Sprintf("/%s/{%s}/{%s}",
		config.GetValueMethod,
		constants.MetTypeURL,
		constants.MetNameURL,
	)
	s.router.Route(fmt.Sprintf("/%s", constants.GetValueMethod), func(r chi.Router) {
		r.Get(str, func(w http.ResponseWriter, r *http.Request) {
			GetMetricValueHandler(w, r, s.store)

		})
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			GetMetricValueJSONHandler(w, r, s.store)

		})

	})
}

package server

import (
	"fmt"
	"net/http"

	"github.com/golovanevvs/metalecoll/internal/server/constants"
)

func (s *server) configureRouter(config *Config) {
	s.router.Use(func(h http.Handler) http.Handler {
		return WithLogging(h)
	})

	str := fmt.Sprintf("/%s/{%s}/{%s}/{%s}",
		config.UpdateMethod,
		constants.MetTypeURL,
		constants.MetNameURL,
		constants.MetValueURL,
	)
	s.router.Post(str, func(w http.ResponseWriter, r *http.Request) {
		UpdateMetricsHandler(w, r, s.store)
	})

	s.router.Get("/", GetMetricNamesHandler)
	str = fmt.Sprintf("/%s/{%s}/{%s}",
		config.GetValueMethod,
		constants.MetTypeURL,
		constants.MetNameURL,
	)
	s.router.Get(str, GetMetricValueHandler)
}

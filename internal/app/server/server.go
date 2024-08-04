package server

import "net/http"

func Start(config *Config) error {
	srv := newServer()
	return http.ListenAndServe(config.bindAddr, srv)
}

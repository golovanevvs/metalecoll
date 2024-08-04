package server

type Config struct {
	bindAddr string
}

func NewConfig() *Config {
	return &Config{
		bindAddr: "localhost:8080",
	}
}

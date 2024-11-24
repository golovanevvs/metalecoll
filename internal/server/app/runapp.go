package app

import (
	"github.com/golovanevvs/metalecoll/internal/server/config"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	// инициализация логгера
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)

	// инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		lg.Fatalf("ошибка инициализации конфигурации сервера: %s", err.Error())
	}

	
}

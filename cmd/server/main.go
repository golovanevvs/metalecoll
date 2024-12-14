// Модуль main обеспечивает запуск сервера.
package main

import (
	"fmt"
	_ "net/http/pprof"

	"github.com/golovanevvs/metalecoll/internal/server/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	app.RunApp()
}

package main

import (
	_ "net/http/pprof"

	"github.com/golovanevvs/metalecoll/internal/server/app"
)

func main() {
	app.RunApp()
}

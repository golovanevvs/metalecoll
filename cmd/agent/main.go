// Модуль main обеспечивает запуск агента.
package main

import (
	"fmt"

	"github.com/golovanevvs/metalecoll/internal/agent/agent"
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
	config, err := agent.NewConfig()
	if err != nil {
		panic(err)
	}
	agent.Start(config)
}

// Модуль main обеспечивает запуск агента.
package main

import (
	"fmt"
	"time"

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
		fmt.Println(err)
		time.Sleep(time.Second * 15)
		return
	}
	agent.Start(config)
}

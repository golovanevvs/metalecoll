// Модуль main обеспечивает запуск агента.
package main

import (
	"github.com/golovanevvs/metalecoll/internal/agent/agent"
)

func main() {
	config := agent.NewConfig()
	agent.Start(config)
}

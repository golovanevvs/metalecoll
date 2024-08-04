package main

import "github.com/golovanevvs/metalecoll/internal/app/server"

func main() {
	config := server.NewConfig()

	if err := server.Start(config); err != nil {
		panic(err)
	}
}

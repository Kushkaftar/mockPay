package main

import (
	"flag"
	"log"
	"mockPay/internal/pkg/app"
	"mockPay/pkg/config"
)

const pathToConfig = "./configs"

func main() {
	var fileName string

	flag.StringVar(&fileName, "env", "", "desc")
	flag.Parse()

	c, err := config.NewConfig(fileName, pathToConfig)
	if err != nil {
		log.Fatalf("failed to load config, err - %s", err)
	}

	if err := app.Start(c); err != nil {
		log.Fatalf("app not start, error - %s", err)
	}
}

package main

import (
	"flag"
	"mockPay/internal/pkg/app"
	"mockPay/pkg/config"
)

const pathToConfig = "./configs"

func main() {
	var fileName string

	flag.StringVar(&fileName, "env", "", "use flag \"-env\" for config file name")
	flag.Parse()

	c := config.MustConfig(fileName, pathToConfig)

	app.MustStart(c)
}

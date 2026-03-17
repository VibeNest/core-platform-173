package main

import (
	"log"

	"gitverse.ru/apavlov-systems/core-platform/config"
	"gitverse.ru/apavlov-systems/core-platform/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}

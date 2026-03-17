package app

import (
	"log"

	"gitverse.ru/apavlov-systems/core-platform/config"
)

// Run — “продолжение main”: тут будет DI и запуск серверов.
func Run(cfg *config.Config) {
	_ = cfg
	if cfg.App.Env == "dev" {
		log.Println("🛠 Запущено в режиме РАЗРАБОТКИ.")
	} else {
		log.Println("🚀 Запущено в режиме ПРОДАКШЕН. Максимальная производительность.")
	}
}

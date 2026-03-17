package app

import (
	"fmt"

	"gitverse.ru/apavlov-systems/core-platform/config"
)

// Run — “продолжение main”: тут будет DI и запуск серверов.
func Run(cfg *config.Config) {
	_ = cfg
	fmt.Println("Проверка выполенения модуля")
}

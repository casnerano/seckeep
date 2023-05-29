package main

import (
	"log"

	"github.com/casnerano/seckeep/internal/client"
)

func main() {
	app, err := client.NewApp()
	if err != nil {
		log.Fatal("Ошибка инициализации приложения.", err.Error())
	}
	defer app.Shutdown()

	if err = app.Run(); err != nil {
		log.Fatal("Ошибка запуска приложения.", err.Error())
	}
}

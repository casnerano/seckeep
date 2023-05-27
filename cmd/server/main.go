package main

import (
	"log"

	"github.com/casnerano/seckeep/internal/server"
)

func main() {
	app, err := server.NewApp()
	if err != nil {
		log.Fatal("Ошибка инициализации приложения.", err.Error())
	}
	defer app.Shutdown()

	if err = app.Run(); err != nil {
		log.Fatal("Ошибка старта приложения.", err.Error())
	}
}

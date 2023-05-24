package main

import (
	"github.com/casnerano/seckeep/internal/client"
)

func main() {
	app, err := client.NewApp()
	if err != nil {
		panic(err.Error())
	}
	defer app.Shutdown()

	app.Run()
}

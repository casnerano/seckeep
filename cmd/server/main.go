package main

import "github.com/casnerano/seckeep/internal/server"

func main() {
	app, err := server.NewApp()
	if err != nil {
		panic(err.Error())
	}
	defer app.Shutdown()

	app.Run()
}

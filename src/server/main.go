package main

import (
	"mofaried/api"
)

func main() {
	app := api.App{
		Port:   ":8000",
		DBPath: "../../ecommerce.db",
	}
	app.Initialize()
	app.Run()
}

package main

import (
	"github.com/m-faried/api"
)

func main() {
	app := api.App{
		Port:   ":8000",
		DBPath: "../../database/ecommerce.db",
	}
	app.Initialize()
	app.Run()
}

package main

import (
	"mofaried/api"
)

func main() {
	app := api.App{
		Port: ":8000",
	}
	app.Initialize("../../ecommerce.db")
	app.Run()
}

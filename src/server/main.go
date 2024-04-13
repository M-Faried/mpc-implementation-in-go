package main

import (
	"mofaried/backend"
)

func main() {
	app := backend.App{
		Port: ":8000",
	}
	app.Initialize()
	app.Run()
}

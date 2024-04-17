module mofaried/server

go 1.22.2

replace mofaried/api => ../api

require mofaried/api v0.0.0-00010101000000-000000000000

require (
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
)

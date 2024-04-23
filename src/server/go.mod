module github.com/m-faried/server

go 1.22.2

replace github.com/m-faried/api => ../api

require github.com/m-faried/api v0.0.0-00010101000000-000000000000

require (
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
)

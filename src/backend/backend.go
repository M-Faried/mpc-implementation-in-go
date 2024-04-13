package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db   *sql.DB
	Port string
}

func helloWorld(res http.ResponseWriter, req *http.Request) {
	log.Println("A new request received!")
	fmt.Fprintf(res, "Hello World")
}

func (a *App) Initialize() {
	database, err := sql.Open("sqlite3", "../../practiceit.db")
	a.db = database
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (a *App) Run() {
	http.HandleFunc("/", helloWorld)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, nil))
}

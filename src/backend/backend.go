package backend

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	database *sql.DB
	router   *mux.Router
	Port     string
}

func helloWorld(res http.ResponseWriter, req *http.Request) {
	log.Println("A new request received!")
	fmt.Fprintf(res, "Hello World")
}

func (a *App) Initialize() {
	database, err := sql.Open("sqlite3", "../../practiceit.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	a.database = database
	a.router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/", helloWorld).Methods("GET")
}

func (a *App) Run() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.router))
}

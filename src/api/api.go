package api

import (
	"fmt"
	"log"
	ctrls "mofaried/api/controllers"
	"mofaried/api/models"
	"mofaried/api/routers"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Port     string
	DBPath   string
	database *models.EcommerceDal
	router   *mux.Router
}

func (a *App) Initialize() {

	if a.DBPath == "" || a.Port == "" {
		panic("App configuration is missing")
	}

	// Creating & initializaing the datasource
	var source models.EcommerceDal
	source.Initialize(a.DBPath)
	a.database = &source

	// Creating & initializaing the router
	a.router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.router))
}

func (a *App) initializeRoutes() {
	a.router.HandleFunc("/", healthCheck).Methods("GET")

	pc := ctrls.NewProductsController(a.database)
	pr := routers.NewProductsRouter(a.router, pc)
	pr.InitRoutes()

	oc := ctrls.NewOrdersController(a.database)
	or := routers.NewOrdersRouter(a.router, oc)
	or.InitRoutes()
}

//////////////////// Helper Functions

func healthCheck(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello World")
}

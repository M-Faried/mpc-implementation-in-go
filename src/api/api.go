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
	DataSource *models.EcommerceDal
	Router     *mux.Router
	Port       string
}

func (a *App) Initialize(dbPath string) {
	// Creating & initializaing the datasource
	var source models.EcommerceDal
	source.Initialize(dbPath)
	a.DataSource = &source
	// Creating & initializaing the router
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", healthCheck).Methods("GET")

	pc := ctrls.NewProductsController(a.DataSource)
	pr := routers.NewProductsRouter(a.Router, pc)
	pr.InitRoutes()

	oc := ctrls.NewOrdersController(a.DataSource)
	or := routers.NewOrdersRouter(a.Router, oc)
	or.InitRoutes()
}

func (a *App) Run() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

//////////////////// Helper Functions

func healthCheck(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello World")
}

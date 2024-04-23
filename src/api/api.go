package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ctrls "github.com/m-faried/api/controllers"
	models "github.com/m-faried/api/models"
	httpRouters "github.com/m-faried/api/presentations/routers"
)

type App struct {
	Port     string
	DBPath   string
	database *models.EcommerceDal
	router   *mux.Router
}

func (a *App) Initialize() {
	// Config validations
	if a.DBPath == "" || a.Port == "" {
		panic("App configuration is missing")
	}

	// Creating & initializaing the datasource
	var source models.EcommerceDal
	source.Initialize(a.DBPath)
	a.database = &source

	// Creating & initializaing the router
	a.initHttpRouter()
}

func (a *App) Run() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.router))
}

func (a *App) initHttpRouter() {

	router := mux.NewRouter()

	router.HandleFunc("/", healthCheck).Methods("GET")

	pc := ctrls.NewProductsController(a.database)
	pr := httpRouters.NewProductsRouter(router, pc)
	pr.InitRoutes()

	oc := ctrls.NewOrdersController(a.database)
	or := httpRouters.NewOrdersRouter(router, oc)
	or.InitRoutes()

	a.router = router
}

//////////////////// Helper Functions

func healthCheck(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello World")
}

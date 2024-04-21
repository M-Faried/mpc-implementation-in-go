package api

import (
	"fmt"
	"log"
	ctrls "mofaried/api/controllers"
	"mofaried/api/handlers"
	"mofaried/api/models"
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
	pr := handlers.NewProductRouter(pc)
	a.Router.HandleFunc("/products", pr.GetAllProducts).Methods("GET")
	a.Router.HandleFunc("/products/{id}", pr.GetSingleProduct).Methods("GET")
	a.Router.HandleFunc("/products", pr.CreateNewProduct).Methods("POST")

	oc := ctrls.NewOrdersController(a.DataSource)
	or := handlers.NewOrdersRoutesHandler(oc)
	a.Router.HandleFunc("/orders", or.GetAllOrders).Methods("GET")
	a.Router.HandleFunc("/orders/{id}", or.GetSingleOrder).Methods("GET")
	a.Router.HandleFunc("/orders", or.CreateNewOrder).Methods("POST")
	a.Router.HandleFunc("/orderitems", or.CreateNewOrderItems).Methods("POST")
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

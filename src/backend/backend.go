package backend

import (
	"fmt"

	"log"
	"mofaried/backend/controllers"
	"mofaried/backend/models"
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

	pc := controllers.ProductsController{
		DS: a.DataSource,
	}
	a.Router.HandleFunc("/products", pc.GetAllProducts).Methods("GET")
	a.Router.HandleFunc("/products/{id}", pc.GetSingleProduct).Methods("GET")
	a.Router.HandleFunc("/products", pc.CreateNewProduct).Methods("POST")

	oc := controllers.OrdersController{
		DS: a.DataSource,
	}
	a.Router.HandleFunc("/orders", oc.GetAllOrders).Methods("GET")
	a.Router.HandleFunc("/orders/{id}", oc.GetSingleOrder).Methods("GET")
	a.Router.HandleFunc("/orders", oc.CreateNewOrder).Methods("POST")
	a.Router.HandleFunc("/orderitems", oc.CreateNewOrderItem).Methods("POST")
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

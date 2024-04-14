package backend

import (
	"database/sql"
	"fmt"

	"log"
	"mofaried/backend/controllers"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	database *sql.DB
	router   *mux.Router
	Port     string
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
	a.router.HandleFunc("/", healthCheck).Methods("GET")

	pc := controllers.ProductsController{
		Database: a.database,
	}
	a.router.HandleFunc("/products", pc.GetAllProducts).Methods("GET")
	a.router.HandleFunc("/products/{id}", pc.GetSingleProduct).Methods("GET")
	a.router.HandleFunc("/products", pc.CreateNewProduct).Methods("POST")

	oc := controllers.OrdersController{
		Database: a.database,
	}
	a.router.HandleFunc("/orders", oc.GetAllOrders).Methods("GET")
	a.router.HandleFunc("/orders/{id}", oc.GetSingleOrder).Methods("GET")
	a.router.HandleFunc("/orders", oc.CreateNewOrder).Methods("POST")
	a.router.HandleFunc("/orderitems", oc.CreateNewOrderItem).Methods("POST")
}

func (a *App) Run() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.router))
}

//////////////////// Helper Functions

func healthCheck(res http.ResponseWriter, req *http.Request) {
	log.Println("A new request received!")
	fmt.Fprintf(res, "Hello World")
}

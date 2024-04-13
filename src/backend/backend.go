package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

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
	a.router.HandleFunc("/products", a.allProducts).Methods("GET")
	a.router.HandleFunc("/products/{id}", a.fetchProduct).Methods("GET")
	a.router.HandleFunc("/products", a.newProduct).Methods("POST")

	a.router.HandleFunc("/orders", a.allOrders).Methods("GET")
	a.router.HandleFunc("/orders/{id}", a.fetchOrder).Methods("GET")
	a.router.HandleFunc("/orders", a.newOrder).Methods("POST")
	a.router.HandleFunc("/orderitems", a.newOrderItem).Methods("POST")
}

func (a *App) Run() {
	// You can do the following as well.
	// http.Handle("/", a.router)
	fmt.Println("Server started and listening on the port", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.router))
}

func (a *App) allProducts(res http.ResponseWriter, req *http.Request) {
	products, err := getProducts(a.database)
	if err != nil {
		fmt.Printf("getProducts err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, products)
}

func (a *App) fetchProduct(res http.ResponseWriter, req *http.Request) {

	// Parsing the submitted id.
	vars := mux.Vars(req)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("fetchProduct invalid product ID: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	// Reading the corresponding product from the database.
	var prod product
	prod.ID = intID
	err2 := prod.getProduct(a.database)
	if err2 != nil {
		fmt.Printf("fetchProduct err: %s\n", err2.Error())
		respondWithError(res, http.StatusNotFound, "Product ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, prod)
}

func (a *App) newProduct(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var p product
	json.Unmarshal(reqBody, &p)
	err := p.newProduct(a.database)

	if err != nil {
		fmt.Printf("newProduct error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(res, http.StatusOK, p)
}

func (a *App) allOrders(res http.ResponseWriter, req *http.Request) {
	orders, err := getAllOrders(a.database)
	if err != nil {
		fmt.Printf("allOrders err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, orders)
}

func (a *App) fetchOrder(res http.ResponseWriter, req *http.Request) {
	// Parsing the submitted id.
	vars := mux.Vars(req)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("fetchOrder invalid order ID: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, "Invalid Order ID")
		return
	}

	// Reading the corresponding product from the database.
	var o order
	o.ID = intID
	err = o.getOrder(a.database)
	if err != nil {
		fmt.Printf("fetchOrder err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Product ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, o)
}

func (a *App) newOrder(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var o order
	json.Unmarshal(reqBody, &o)
	err := o.createOrder(a.database)
	if err != nil {
		fmt.Printf("newOrder error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	for _, item := range o.Items {
		item.OrderID = o.ID
		err := item.createOrderItem(a.database)
		if err != nil {
			fmt.Printf("newOrder error: %s\n", err.Error())
			respondWithError(res, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(res, http.StatusOK, o)
}

func (a *App) newOrderItem(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var items []orderItem
	json.Unmarshal(reqBody, &items)
	for _, item := range items {
		err := item.createOrderItem(a.database)
		if err != nil {
			fmt.Printf("newOrderItem error: %s\n", err.Error())
			respondWithError(res, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(res, http.StatusOK, items)
}

//////////////////// Helper Functions

func healthCheck(res http.ResponseWriter, req *http.Request) {
	log.Println("A new request received!")
	fmt.Fprintf(res, "Hello World")
}

// respondWithError sets the error response to the response writer
func respondWithError(res http.ResponseWriter, code int, message string) {
	respondWithJSON(res, code, map[string]string{"error": message})
}

// respondWithJSON sets the response payload in the response writer
func respondWithJSON(res http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(response)
}

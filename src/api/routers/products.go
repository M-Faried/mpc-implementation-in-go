package routers

import (
	"encoding/json"
	"fmt"
	"io"
	"mofaried/api/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type IProductsController interface {
	GetProducts() ([]models.Product, error)
	GetSingleProduct(int) (*models.Product, error)
	CreateProduct(*models.Product) error
}

type productsRouter struct {
	ctrl   IProductsController
	router *mux.Router
}

func NewProductsRouter(router *mux.Router, controller IProductsController) IRouter {
	return &productsRouter{
		ctrl:   controller,
		router: router,
	}
}

func (pr *productsRouter) InitRoutes() {
	pr.router.HandleFunc("/products", pr.getAllProductsHandler).Methods("GET")
	pr.router.HandleFunc("/products/{id}", pr.getSingleProductHandler).Methods("GET")
	pr.router.HandleFunc("/products", pr.createNewProductHandler).Methods("POST")
}

func (pr *productsRouter) getAllProductsHandler(res http.ResponseWriter, req *http.Request) {
	products, err := pr.ctrl.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, products)
}

func (pr *productsRouter) getSingleProductHandler(res http.ResponseWriter, req *http.Request) {

	// Parsing the submitted id.
	vars := mux.Vars(req)
	paramID := vars["id"]
	id, err := strconv.Atoi(paramID)
	if err != nil {
		fmt.Printf("GetSingleProduct invalid product ID: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p, err := pr.ctrl.GetSingleProduct(id)
	if err != nil {
		fmt.Printf("GetSingleProduct err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Product ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, p)
}

func (pr *productsRouter) createNewProductHandler(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var p models.Product
	json.Unmarshal(reqBody, &p)
	err := pr.ctrl.CreateProduct(&p)

	if err != nil {
		fmt.Printf("CreateNewProduct error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(res, http.StatusOK, p)
}

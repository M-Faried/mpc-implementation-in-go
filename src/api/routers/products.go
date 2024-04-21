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

type ProductsRouters struct {
	ctrl   IProductsController
	router *mux.Router
}

func NewProductsRouter(router *mux.Router, controller IProductsController) *ProductsRouters {
	return &ProductsRouters{
		ctrl:   controller,
		router: router,
	}
}

func (pr *ProductsRouters) InitRoutes() {
	pr.router.HandleFunc("/products", pr.getAllProducts).Methods("GET")
	pr.router.HandleFunc("/products/{id}", pr.getSingleProduct).Methods("GET")
	pr.router.HandleFunc("/products", pr.createNewProduct).Methods("POST")
}

func (pr *ProductsRouters) getAllProducts(res http.ResponseWriter, req *http.Request) {
	products, err := pr.ctrl.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, products)
}

func (pr *ProductsRouters) getSingleProduct(res http.ResponseWriter, req *http.Request) {

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

func (pr *ProductsRouters) createNewProduct(res http.ResponseWriter, req *http.Request) {
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

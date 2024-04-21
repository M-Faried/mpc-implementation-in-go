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

type IProductController interface {
	GetProducts() (*[]models.Product, error)
	FindProductByID(int) (*models.Product, error)
	CreateProduct(*models.Product) error
}

type ProductRouter struct {
	ctrl IProductController
}

func NewProductRouter(controller IProductController) *ProductRouter {
	return &ProductRouter{
		ctrl: controller,
	}
}

func (pr *ProductRouter) GetAllProducts(res http.ResponseWriter, req *http.Request) {
	products, err := pr.ctrl.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, products)
}

func (pr *ProductRouter) GetSingleProduct(res http.ResponseWriter, req *http.Request) {

	// Parsing the submitted id.
	vars := mux.Vars(req)
	paramID := vars["id"]
	id, err := strconv.Atoi(paramID)
	if err != nil {
		fmt.Printf("GetSingleProduct invalid product ID: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	// Reading the corresponding product from the database.
	// p := models.Product{
	// 	ID: id,
	// }
	p, err := pr.ctrl.FindProductByID(id)
	if err != nil {
		fmt.Printf("GetSingleProduct err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Product ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, p)
}

func (pr *ProductRouter) CreateNewProduct(res http.ResponseWriter, req *http.Request) {
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

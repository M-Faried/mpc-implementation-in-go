package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mofaried/api/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type IProductsDataSource interface {
	GetProducts() ([]models.Product, error)
	FindProductByID(p *models.Product) error
	CreateProduct(p *models.Product) error
}

type ProductsController struct {
	ds IProductsDataSource
}

func NewProductsController(dataSource IProductsDataSource) *ProductsController {
	return &ProductsController{
		ds: dataSource,
	}
}

func (c *ProductsController) GetAllProducts(res http.ResponseWriter, req *http.Request) {
	products, err := c.ds.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, products)
}

func (c *ProductsController) GetSingleProduct(res http.ResponseWriter, req *http.Request) {

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
	p := models.Product{
		ID: id,
	}
	err = c.ds.FindProductByID(&p)
	if err != nil {
		fmt.Printf("GetSingleProduct err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Product ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, p)
}

func (c *ProductsController) CreateNewProduct(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var p models.Product
	json.Unmarshal(reqBody, &p)
	err := c.ds.CreateProduct(&p)

	if err != nil {
		fmt.Printf("CreateNewProduct error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(res, http.StatusOK, p)
}

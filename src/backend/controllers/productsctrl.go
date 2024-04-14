package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"mofaried/backend/models"

	"github.com/gorilla/mux"
)

type ProductsController struct {
	Database *sql.DB
}

func (c *ProductsController) GetAllProducts(res http.ResponseWriter, req *http.Request) {
	products, err := c.dbGetProducts()
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
	err = c.dbFindProductByID(&p)
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
	err := c.dbCreateProduct(&p)

	if err != nil {
		fmt.Printf("CreateNewProduct error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(res, http.StatusOK, p)
}

//////////////// DB functions

func (c *ProductsController) dbGetProducts() ([]models.Product, error) {
	rows, err := c.Database.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []models.Product{}
	var p models.Product

	for rows.Next() {
		err := rows.Scan(&p.ID, &p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status)
		if err != nil {
			return products, err
		} else {
			products = append(products, p)
		}
	}

	return products, nil
}

func (c *ProductsController) dbFindProductByID(p *models.Product) error {
	query := `
		SELECT productCode, name, inventory, price, status 
		FROM products 
		WHERE id=?
	`
	row := c.Database.QueryRow(query, p.ID)

	err := row.Scan(&p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status)
	return err
}

func (c *ProductsController) dbCreateProduct(p *models.Product) error {
	query := `
		INSERT INTO products(productCode, name, inventory, price, status) 
		VALUES(?, ?, ?, ?, ?)
	`
	res, err := c.Database.Exec(query, p.ProductCode, p.Name, p.Inventory, p.Price, p.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

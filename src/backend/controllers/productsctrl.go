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
	products, err := dbGetProducts(c.Database)
	if err != nil {
		fmt.Printf("getProducts err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, products)
}

func (c *ProductsController) GetSingleProduct(res http.ResponseWriter, req *http.Request) {

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
	prod, err := dbFindProductByID(c.Database, intID)
	if err != nil {
		fmt.Printf("fetchProduct err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Product ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, prod)
}

func (c *ProductsController) CreateNewProduct(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var p models.Product
	json.Unmarshal(reqBody, &p)
	err := dbCreateProduct(c.Database, &p)

	if err != nil {
		fmt.Printf("newProduct error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(res, http.StatusOK, p)
}

//////////////// DB functions

func dbGetProducts(db *sql.DB) ([]models.Product, error) {
	rows, err := db.Query("SELECT * FROM products")
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

func dbFindProductByID(db *sql.DB, prodID int) (models.Product, error) {
	query := `
		SELECT productCode, name, inventory, price, status 
		FROM products 
		WHERE id=?
	`
	row := db.QueryRow(query, prodID)
	p := models.Product{
		ID: prodID,
	}
	err := row.Scan(&p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status)
	return p, err
}

func dbCreateProduct(db *sql.DB, p *models.Product) error {
	query := `
		INSERT INTO products(productCode, name, inventory, price, status) 
		VALUES(?, ?, ?, ?, ?)
	`
	res, err := db.Exec(query, p.ProductCode, p.Name, p.Inventory, p.Price, p.Status)
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

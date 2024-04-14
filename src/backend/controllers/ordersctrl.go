package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mofaried/backend/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrdersController struct {
	Database *sql.DB
}

func (c *OrdersController) GetAllOrders(res http.ResponseWriter, req *http.Request) {
	orders, err := c.dbGetAllOrders()
	if err != nil {
		fmt.Printf("GetAllOrders err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, orders)
}

func (c *OrdersController) GetSingleOrder(res http.ResponseWriter, req *http.Request) {
	// Parsing the submitted id.
	vars := mux.Vars(req)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("GetSingleOrder invalid order ID: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, "Invalid Order ID")
		return
	}

	// Reading the corresponding order from the database.
	var o models.Order
	o.ID = intID
	err = c.dbGetOrderByID(&o)
	if err != nil {
		fmt.Printf("GetSingleOrder err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Order ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, o)
}

func (c *OrdersController) CreateNewOrder(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var o models.Order
	json.Unmarshal(reqBody, &o)
	err := c.dbCreateOrder(&o)
	if err != nil {
		fmt.Printf("CreateNewOrder error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	for _, item := range o.Items {
		item.OrderID = o.ID
		err := c.dbCreateOrderItem(&item)
		if err != nil {
			fmt.Printf("CreateNewOrder error: %s\n", err.Error())
			respondWithError(res, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(res, http.StatusOK, o)
}

func (c *OrdersController) CreateNewOrderItem(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var items []models.OrderItem
	json.Unmarshal(reqBody, &items)
	for _, item := range items {
		err := c.dbCreateOrderItem(&item)
		if err != nil {
			fmt.Printf("CreateNewOrderItem error: %s\n", err.Error())
			respondWithError(res, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(res, http.StatusOK, items)
}

/////////////// DB Functions

func (c *OrdersController) dbGetAllOrders() ([]models.Order, error) {
	rows, err := c.Database.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		err := rows.Scan(&o.ID, &o.CustomerName, &o.Total, &o.Status)
		if err != nil {
			return orders, err
		} else {
			err = c.dbGetOrderItems(&o)
			if err != nil {
				return orders, err
			}
			orders = append(orders, o)
		}
	}

	return orders, nil
}

func (c *OrdersController) dbGetOrderItems(o *models.Order) error {
	query := `
		SELECT * 
		FROM order_items
		WHERE order_id=?
	`
	rows, err := c.Database.Query(query, o.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var orderItems []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.OrderID, &item.ProductID, &item.Quantity)
		if err != nil {
			return err
		} else {
			orderItems = append(orderItems, item)
		}
	}
	o.Items = orderItems
	return nil
}

func (c *OrdersController) dbGetOrderByID(o *models.Order) error {
	query := `
		SELECT customerName, total, status
		FROM orders
		WHERE id=?
	`
	row := c.Database.QueryRow(query, o.ID)
	err := row.Scan(&o.CustomerName, &o.Total, &o.Status)
	if err != nil {
		return err
	}
	err = c.dbGetOrderItems(o)
	if err != nil {
		return err
	}
	return nil
}

func (c *OrdersController) dbCreateOrder(o *models.Order) error {
	query := `
		INSERT INTO orders(customerName, total, status) 
		VALUES(?, ?, ?)
	`
	result, err := c.Database.Exec(query, o.CustomerName, o.Total, o.Status)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	o.ID = int(id)
	return nil
}

func (c *OrdersController) dbCreateOrderItem(item *models.OrderItem) error {
	query := `
		INSERT INTO order_items(order_id, product_id, quantity) 
		VALUES(?, ?, ?)
	`
	_, err := c.Database.Exec(query, item.OrderID, item.ProductID, item.Quantity)
	if err != nil {
		return err
	}
	return nil
}

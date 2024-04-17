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

type IOrdersDataSource interface {
	GetAllOrders() ([]models.Order, error)
	GetOrderItems(o *models.Order) error
	GetOrderByID(o *models.Order) error
	CreateOrder(o *models.Order) error
	CreateOrderItem(item *models.OrderItem) error
}

type OrdersController struct {
	DS IOrdersDataSource
}

func (c *OrdersController) GetAllOrders(res http.ResponseWriter, req *http.Request) {
	orders, err := c.DS.GetAllOrders()
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
	err = c.DS.GetOrderByID(&o)
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
	err := c.DS.CreateOrder(&o)
	if err != nil {
		fmt.Printf("CreateNewOrder error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	for _, item := range o.Items {
		item.OrderID = o.ID
		err := c.DS.CreateOrderItem(&item)
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
		err := c.DS.CreateOrderItem(&item)
		if err != nil {
			fmt.Printf("CreateNewOrderItem error: %s\n", err.Error())
			respondWithError(res, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondWithJSON(res, http.StatusOK, items)
}

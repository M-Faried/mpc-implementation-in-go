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

type IOrdersController interface {
	GetAllOrders() ([]models.Order, error)
	GetOrderItems(o *models.Order) ([]models.OrderItem, error)
	GetOrderByID(int) (*models.Order, error)
	CreateOrder(o *models.Order) error
	AddOrderItems(items []models.OrderItem) error
}

type OrdersRouter struct {
	ctrl IOrdersController
}

func NewOrdersRouter(controller IOrdersController) *OrdersRouter {
	return &OrdersRouter{
		ctrl: controller,
	}
}

func (or *OrdersRouter) GetAllOrders(res http.ResponseWriter, req *http.Request) {
	orders, err := or.ctrl.GetAllOrders()
	if err != nil {
		fmt.Printf("GetAllOrders err: %s\n", err.Error())
		respondWithError(res, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, orders)
}

func (or *OrdersRouter) GetSingleOrder(res http.ResponseWriter, req *http.Request) {
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
	o, err := or.ctrl.GetOrderByID(intID)
	if err != nil {
		fmt.Printf("GetSingleOrder err: %s\n", err.Error())
		respondWithError(res, http.StatusNotFound, "Order ID Is Not Found")
		return
	}
	respondWithJSON(res, http.StatusOK, o)
}

func (or *OrdersRouter) CreateNewOrder(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var o models.Order
	json.Unmarshal(reqBody, &o)
	err := or.ctrl.CreateOrder(&o)
	if err != nil {
		fmt.Printf("CreateNewOrder error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}

	// for _, item := range o.Items {
	// 	item.OrderID = o.ID
	// 	err := or.ctrl.CreateOrderItem(&item)
	// 	if err != nil {
	// 		fmt.Printf("CreateNewOrder error: %s\n", err.Error())
	// 		respondWithError(res, http.StatusBadRequest, err.Error())
	// 		return
	// 	}
	// }

	respondWithJSON(res, http.StatusOK, o)
}

func (or *OrdersRouter) AddOrderItems(res http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var items []models.OrderItem
	json.Unmarshal(reqBody, &items)
	err := or.ctrl.AddOrderItems(items)
	if err != nil {
		fmt.Printf("AddOrderItems error: %s\n", err.Error())
		respondWithError(res, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(res, http.StatusOK, items)
}
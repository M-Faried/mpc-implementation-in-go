package controllers

import (
	"github.com/m-faried/api/models"
)

type IOrdersDataSource interface {
	GetAllOrders() ([]models.Order, error)
	GetOrderItems(o *models.Order) error
	GetOrderByID(o *models.Order) error
	CreateOrder(o *models.Order) error
	CreateOrderItem(item *models.OrderItem) error
}

type OrdersController struct {
	ds IOrdersDataSource
}

func NewOrdersController(dataSource IOrdersDataSource) *OrdersController {
	return &OrdersController{
		ds: dataSource,
	}
}

func (c *OrdersController) GetAllOrders() ([]models.Order, error) {
	orders, err := c.ds.GetAllOrders()
	return orders, err
}

func (c *OrdersController) GetOrderByID(id int) (*models.Order, error) {
	// Reading the corresponding order from the database.
	var o models.Order
	o.ID = id
	err := c.ds.GetOrderByID(&o)
	return &o, err
}

func (c *OrdersController) CreateOrder(o *models.Order) error {
	// Creating the order.
	err := c.ds.CreateOrder(o)
	if err != nil {
		return err
	}

	// Creating the order items.
	for _, item := range o.Items {
		item.OrderID = o.ID
		err := c.ds.CreateOrderItem(&item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *OrdersController) AddOrderItems(items []models.OrderItem) error {
	for _, item := range items {
		err := c.ds.CreateOrderItem(&item)
		if err != nil {
			return err
		}
	}
	return nil
}

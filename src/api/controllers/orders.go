package controllers

import (
	t "github.com/m-faried/types"
)

type IOrdersDataSource interface {
	GetAllOrders() ([]t.Order, error)
	GetOrderItems(o *t.Order) error
	GetOrderByID(o *t.Order) error
	CreateOrder(o *t.Order) error
	CreateOrderItem(item *t.OrderItem) error
}

type OrdersController struct {
	ds IOrdersDataSource
}

func NewOrdersController(dataSource IOrdersDataSource) *OrdersController {
	return &OrdersController{
		ds: dataSource,
	}
}

func (c *OrdersController) GetAllOrders() ([]t.Order, error) {
	orders, err := c.ds.GetAllOrders()
	return orders, err
}

func (c *OrdersController) GetOrderByID(id int) (*t.Order, error) {
	// Reading the corresponding order from the database.
	var o t.Order
	o.ID = id
	err := c.ds.GetOrderByID(&o)
	return &o, err
}

func (c *OrdersController) CreateOrder(o *t.Order) error {
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

func (c *OrdersController) AddOrderItems(items []t.OrderItem) error {
	for _, item := range items {
		err := c.ds.CreateOrderItem(&item)
		if err != nil {
			return err
		}
	}
	return nil
}

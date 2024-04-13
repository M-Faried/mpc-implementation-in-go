package backend

import (
	"database/sql"
)

type order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Total        int         `json:"total"`
	Status       string      `json:"status"`
	Items        []orderItem `json:"items"`
}

type orderItem struct {
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func getAllOrders(db *sql.DB) ([]order, error) {
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []order{}
	for rows.Next() {
		var o order
		err := rows.Scan(&o.ID, &o.CustomerName, &o.Total, &o.Status)
		if err != nil {
			return orders, err
		} else {
			err = o.getOrderItems(db)
			if err != nil {
				return orders, err
			}
			orders = append(orders, o)
		}
	}

	return orders, nil
}

func (o *order) getOrderItems(db *sql.DB) error {
	query := `
		SELECT * 
		FROM order_items
		WHERE order_id=?
	`
	rows, err := db.Query(query, o.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	var orderItems []orderItem
	for rows.Next() {
		var item orderItem
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

func (o *order) getOrder(db *sql.DB) error {
	query := `
		SELECT customerName, total, status
		FROM orders
		WHERE id=?
	`
	row := db.QueryRow(query, o.ID)
	err := row.Scan(&o.CustomerName, &o.Total, &o.Status)
	if err != nil {
		return err
	}
	err = o.getOrderItems(db)
	if err != nil {
		return err
	}
	return nil
}

func (o *order) createOrder(db *sql.DB) error {
	query := `
		INSERT INTO orders(customerName, total, status) 
		VALUES(?, ?, ?)
	`
	result, err := db.Exec(query, o.CustomerName, o.Total, o.Status)
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

func (item *orderItem) createOrderItem(db *sql.DB) error {
	query := `
		INSERT INTO order_items(order_id, product_id, quantity) 
		VALUES(?, ?, ?)
	`
	_, err := db.Exec(query, item.OrderID, item.ProductID, item.Quantity)
	if err != nil {
		return err
	}
	return nil
}

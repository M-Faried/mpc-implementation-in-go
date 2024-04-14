package datasources

import (
	"database/sql"
	"log"
	"mofaried/backend/models"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type EcommerceDataSource struct {
	DB *sql.DB
}

var once sync.Once

func (c *EcommerceDataSource) Initialize() {
	once.Do(func() {
		database, err := sql.Open("sqlite3", "../../ecommerce.db")
		if err != nil {
			log.Fatal(err.Error())
		}
		c.DB = database
	})
}

// Product Actions

func (c EcommerceDataSource) GetProducts() ([]models.Product, error) {
	rows, err := c.DB.Query("SELECT * FROM products")
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

func (c EcommerceDataSource) FindProductByID(p *models.Product) error {
	query := `
		SELECT productCode, name, inventory, price, status 
		FROM products 
		WHERE id=?
	`
	row := c.DB.QueryRow(query, p.ID)

	err := row.Scan(&p.ProductCode, &p.Name, &p.Inventory, &p.Price, &p.Status)
	return err
}

func (c EcommerceDataSource) CreateProduct(p *models.Product) error {
	query := `
		INSERT INTO products(productCode, name, inventory, price, status) 
		VALUES(?, ?, ?, ?, ?)
	`
	res, err := c.DB.Exec(query, p.ProductCode, p.Name, p.Inventory, p.Price, p.Status)
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

// Order Actions

func (c EcommerceDataSource) GetAllOrders() ([]models.Order, error) {
	rows, err := c.DB.Query("SELECT * FROM orders")
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
			err = c.GetOrderItems(&o)
			if err != nil {
				return orders, err
			}
			orders = append(orders, o)
		}
	}

	return orders, nil
}

func (c EcommerceDataSource) GetOrderItems(o *models.Order) error {
	query := `
		SELECT * 
		FROM order_items
		WHERE order_id=?
	`
	rows, err := c.DB.Query(query, o.ID)
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

func (c EcommerceDataSource) GetOrderByID(o *models.Order) error {
	query := `
		SELECT customerName, total, status
		FROM orders
		WHERE id=?
	`
	row := c.DB.QueryRow(query, o.ID)
	err := row.Scan(&o.CustomerName, &o.Total, &o.Status)
	if err != nil {
		return err
	}
	err = c.GetOrderItems(o)
	if err != nil {
		return err
	}
	return nil
}

func (c EcommerceDataSource) CreateOrder(o *models.Order) error {
	query := `
		INSERT INTO orders(customerName, total, status) 
		VALUES(?, ?, ?)
	`
	result, err := c.DB.Exec(query, o.CustomerName, o.Total, o.Status)
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

func (c EcommerceDataSource) CreateOrderItem(item *models.OrderItem) error {
	query := `
		INSERT INTO order_items(order_id, product_id, quantity) 
		VALUES(?, ?, ?)
	`
	_, err := c.DB.Exec(query, item.OrderID, item.ProductID, item.Quantity)
	if err != nil {
		return err
	}
	return nil
}

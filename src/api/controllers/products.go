package controllers

import (
	t "github.com/m-faried/types"
)

type IProductsDataSource interface {
	GetProducts() ([]t.Product, error)
	FindProductByID(p *t.Product) error
	CreateProduct(p *t.Product) error
}

type ProductsController struct {
	ds IProductsDataSource
}

func NewProductsController(dataSource IProductsDataSource) *ProductsController {
	return &ProductsController{
		ds: dataSource,
	}
}

func (c *ProductsController) GetProducts() ([]t.Product, error) {
	products, err := c.ds.GetProducts()
	return products, err
}

func (c *ProductsController) GetSingleProduct(id int) (*t.Product, error) {
	// Reading the corresponding product from the database.
	p := t.Product{
		ID: id,
	}
	err := c.ds.FindProductByID(&p)
	return &p, err
}

func (c *ProductsController) CreateProduct(p *t.Product) error {
	err := c.ds.CreateProduct(p)
	return err
}

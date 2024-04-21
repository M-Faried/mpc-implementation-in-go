package controllers

import (
	"mofaried/api/models"
)

type IProductsDataSource interface {
	GetProducts() ([]models.Product, error)
	FindProductByID(p *models.Product) error
	CreateProduct(p *models.Product) error
}

type ProductsController struct {
	ds IProductsDataSource
}

func NewProductsController(dataSource IProductsDataSource) *ProductsController {
	return &ProductsController{
		ds: dataSource,
	}
}

func (c *ProductsController) GetProducts() ([]models.Product, error) {
	products, err := c.ds.GetProducts()
	return products, err
}

func (c *ProductsController) GetSingleProduct(id int) (*models.Product, error) {
	// Reading the corresponding product from the database.
	p := models.Product{
		ID: id,
	}
	err := c.ds.FindProductByID(&p)
	return &p, err
}

func (c *ProductsController) CreateProduct(p *models.Product) error {
	err := c.ds.CreateProduct(p)
	return err
}

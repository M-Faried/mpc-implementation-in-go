package rpc

import (
	"context"
	"fmt"

	"github.com/m-faried/api/models"
)

type IProductsController interface {
	GetProducts() ([]models.Product, error)
}

type IOrdersController interface{}

type GrpcEcommerceServer struct {
	UnimplementedECommerceServer
	prodCtrl  IProductsController
	orderCtrl IOrdersController
}

func NewGrpcEcommerceServer(prodCtrl IProductsController, ordersCtrl IOrdersController) *GrpcEcommerceServer {
	return &GrpcEcommerceServer{
		prodCtrl:  prodCtrl,
		orderCtrl: ordersCtrl,
	}
}

func (s *GrpcEcommerceServer) GetProducts(context.Context, *GetProductsRequest) (*GetProductsResponse, error) {

	products, err := s.prodCtrl.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		return nil, err
	}

	rpcProducts := make([]*Product, 0)

	for _, prod := range products {
		rpcProducts = append(rpcProducts, &Product{
			ID:          int32(prod.ID),
			Price:       int32(prod.Price),
			Inventory:   int32(prod.Inventory),
			Name:        prod.Name,
			ProductCode: prod.ProductCode,
			Status:      prod.Status,
		})
	}

	return &GetProductsResponse{
		Products: rpcProducts,
	}, nil
}

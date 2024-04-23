package rpc

import (
	"context"
	"fmt"

	"github.com/m-faried/api/models"
	rpc "github.com/m-faried/api/presentations/rpc/ecommerce"
	"google.golang.org/grpc"
)

type IProductsController interface {
	GetProducts() ([]models.Product, error)
}

type ProductsServer struct {
	rpc.UnimplementedECommerceServer
	ctrl IProductsController
}

func (s *ProductsServer) GetProducts(ctx context.Context, in *rpc.GetProductsRequest, opts ...grpc.CallOption) (*rpc.GetProductsResponse, error) {

	products, err := s.ctrl.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		return nil, err
	}

	rpcProducts := make([]*rpc.Product, len(products))

	for _, prod := range products {
		rpcProducts = append(rpcProducts, &rpc.Product{
			ID:          int32(prod.ID),
			Price:       int32(prod.Price),
			Inventory:   int32(prod.Inventory),
			Name:        prod.Name,
			ProductCode: prod.ProductCode,
			Status:      prod.Status,
		})
	}

	return &rpc.GetProductsResponse{
		Products: rpcProducts,
	}, nil
}

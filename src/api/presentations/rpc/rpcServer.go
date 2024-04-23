package rpc

import (
	"context"
	"fmt"

	"github.com/m-faried/api/models"
)

// IProductsController is the requried interface for the product controller.
// It contains the methods used by the GrpcEcommerceServer
type IProductsController interface {
	GetProducts() ([]models.Product, error)
	GetSingleProduct(int) (*models.Product, error)
}

// IOrdersController is the requried interface for the order controller
// It contains the methods used by the GrpcEcommerceServer
type IOrdersController interface{}

// GrpcEcommerceServer is the server structure for the gRPC.
// The structure is implementing ECommerceServer interface in
// the generated gRPC go files.
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

	// Reading products data from the database
	products, err := s.prodCtrl.GetProducts()
	if err != nil {
		fmt.Printf("GetAllProducts err: %s\n", err.Error())
		return nil, err
	}

	// Converting models.Product to rpc.Product
	var rpcProducts []*Product
	var currRpcProd *Product
	for _, prod := range products {
		currRpcProd = getRpcProductFromModelProduct(&prod)
		rpcProducts = append(rpcProducts, currRpcProd)
	}

	// Returning the result
	return &GetProductsResponse{
		Products: rpcProducts,
	}, nil
}

func (s *GrpcEcommerceServer) GetSingleProduct(c context.Context, req *GetSignleProductRequest) (*GetSingleProductResponse, error) {

	// Searching for the product into the database
	prodID := int(req.ProductID)
	p, err := s.prodCtrl.GetSingleProduct(prodID)
	if err != nil {
		fmt.Printf("GetSingleProduct err: %s\n", err.Error())
		return nil, err
	}

	// Converting models.Product to rpc.Product
	prod := getRpcProductFromModelProduct(p)

	// Returning the result
	return &GetSingleProductResponse{
		Product: prod,
	}, nil
}

func getRpcProductFromModelProduct(mProd *models.Product) *Product {
	return &Product{
		ID:          int32(mProd.ID),
		Price:       int32(mProd.Price),
		Inventory:   int32(mProd.Inventory),
		Name:        mProd.Name,
		ProductCode: mProd.ProductCode,
		Status:      mProd.Status,
	}
}

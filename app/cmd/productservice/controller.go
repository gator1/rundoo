package main

import (
	"context"
	"log"

	"app/products"
	rundoogrpc "app/api/v1"
)

// productsServiceController implements the gRPC ProductsServiceServer interface.
type productsServiceController struct {
	productsInterface products.ServiceInterface
	rundoogrpc.UnimplementedProductServiceServer
}

// NewProductServiceController instantiates a new ProducterviceServer.
func NewProductsServiceController(productInterface products.ServiceInterface) rundoogrpc.ProductServiceServer {
	return &productsServiceController{
		productsInterface: productInterface,
	}
}

// GetProducts calls the product service's GetProducts method and maps the result to a grpc service response.
func (ctlr *productsServiceController) GetProducts(ctx context.Context, req *rundoogrpc.GetProductsRequest) (resp *rundoogrpc.GetProductsResponse, err error) {
	resultMap, err := ctlr.productsInterface.GetProducts()
	if err != nil {
		return
	}

	resp = &rundoogrpc.GetProductsResponse{}
	for _, u := range resultMap {
		resp.Products = append(resp.Products, marshalProduct(&u))
	}

	log.Printf("handled GetProducts")
	return
}

// marshalProductmarshals a business object Product into a gRPC layer Product.
func marshalProduct(p *products.Product) *rundoogrpc.Product {
	return &rundoogrpc.Product{Sku: string(p.Sku), Category: string(p.Category), Name: p.Name}
}
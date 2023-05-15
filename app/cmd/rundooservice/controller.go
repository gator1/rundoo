package main

import (
	"context"
	"log"

	"app/rundoo"
	rundoogrpc "app/api/v1"
	"app/internal/data"

)

// productsServiceController implements the gRPC ProductsServiceServer interface.
type productsServiceController struct {
	productsInterface rundoo.ServiceInterface
	rundoogrpc.UnimplementedProductServiceServer
}

// NewProductServiceController instantiates a new ProducterviceServer.
func NewProductsServiceController(productInterface rundoo.ServiceInterface) rundoogrpc.ProductServiceServer {
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

	log.Printf("Grpc handled GetProducts")
	return
}

// marshalProductmarshals a business object Product into a gRPC layer Product.
func marshalProduct(p *data.Product) *rundoogrpc.Product {
	return &rundoogrpc.Product{Sku: string(p.Sku), Category: string(p.Category), Name: p.Name}
}

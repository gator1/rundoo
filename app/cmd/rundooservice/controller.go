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
		log.Printf("productsServiceController GetProducts failed %v", err)
		return
	}

	resp = &rundoogrpc.GetProductsResponse{}
	for _, u := range resultMap {
		resp.Products = append(resp.Products, marshalProduct(&u))
	}
	log.Printf("Grpc handled GetProducts")
	return
}

// GetProducts calls the product service's GetProduct method and maps the result to a grpc service response.
func (ctlr *productsServiceController) GetProduct(ctx context.Context, req *rundoogrpc.GetProductRequest) (resp *rundoogrpc.GetProductResponse, err error) {

	result, err := ctlr.productsInterface.GetProduct(req.Id)
	if err != nil {
		log.Printf("productsServiceController GetProduct failed %v", err)
		
		return
	}

	resp = &rundoogrpc.GetProductResponse{}
	resp.Product = marshalProduct(&result)
	log.Printf("Grpc handled GetProduct")
	return
}

// AddProduct calls the product service's AddProduct method and maps the result to a grpc service response.
func (ctlr *productsServiceController) AddProduct(ctx context.Context, req *rundoogrpc.AddProductRequest) (resp *rundoogrpc.AddProductResponse, err error) {

	product := data.Product{
		Name: req.Product.Name,
		Category: data.CategoryType(req.Product.Category),
		Sku: data.SKU(req.Product.Sku),
	}
	result, err := ctlr.productsInterface.AddProduct(product)
	if err != nil {
		log.Printf("productsServiceController AddProduct failed %v", err)
		
		return
	}

	resp = &rundoogrpc.AddProductResponse{}
	resp.Ok = result
	log.Printf("Grpc handled AddProduct")
	return
}


// marshalProductmarshals a business object Product into a gRPC layer Product.
func marshalProduct(p *data.Product) *rundoogrpc.Product {
	return &rundoogrpc.Product{Id: p.ID, Name: p.Name, Category: string(p.Category), Sku: string(p.Sku)}
}

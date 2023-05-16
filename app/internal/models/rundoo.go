package models

import (
	"context"
	"log"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	rundoogrpc "app/api/v1"
	"app/internal/data"
	"app/registry"
)

type Product struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Sku      string   `json:"sku"`
}

type ProductResponse struct {
	Product *Product `json:"product"`
}

type ProductsResponse struct {
	products *[]Product `json:"products"`
}

type RundooModel struct {
	Endpoint string
}

func (m *RundooModel) GetAll() (products []data.Product, err error) {
	log.Println("RundooModel GetAll called!")

	serviceURL, err := registry.GetProvider(registry.RundooService)
	if err != nil {
		log.Println("Error getting provider ProductService: ", err)
		return nil, err
	}
	record := strings.Split(serviceURL, ":") // http://localhost:port
	portInt, _ := strconv.Atoi(record[2])
	rpcPort := ":" + strconv.Itoa(portInt+1)

	conn, err := grpc.Dial("localhost"+rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := rundoogrpc.NewProductServiceClient(conn)

	response, err := client.GetProducts(context.Background(), &rundoogrpc.GetProductsRequest{})
	if err != nil {
		log.Fatalf("failed to get products: %v", err)
		return nil, err
	}

	for _, product := range response.GetProducts() {
		dataproduct :=  data.Product{
			ID: product.Id,
			Name: product.Name,
			Category: data.CategoryType(product.Category),
			Sku: data.SKU(product.Sku),
		}
		products = append(products, dataproduct)
	}
	return 
}

func (m *RundooModel) Get(id int64) (product data.Product, err error) {
	log.Println("RundooModel Get %d called!", id)

	serviceURL, err := registry.GetProvider(registry.RundooService)
	if err != nil {
		log.Println("Error getting provider ProductService: ", err)
		return 
	}
	record := strings.Split(serviceURL, ":") // http://localhost:port
	portInt, _ := strconv.Atoi(record[2])
	rpcPort := ":" + strconv.Itoa(portInt+1)

	conn, err := grpc.Dial("localhost"+rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return 
	}
	defer conn.Close()

	client := rundoogrpc.NewProductServiceClient(conn)

	response, err := client.GetProduct(context.Background(), &rundoogrpc.GetProductRequest{Id: id})
	if err != nil {
		log.Fatalf("failed to get product %v: %v", id, err)
		return 
	}

	gproduct := response.GetProduct()

	product.ID =  gproduct.Id
	product.Name =  gproduct.Name
	product.Category =  data.CategoryType(gproduct.Category)
	product.Sku =  data.SKU(gproduct.Sku)
	return 
}

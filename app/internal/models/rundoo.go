package models

import (
	"context"
	"fmt"
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
	fmt.Println("RundooModel GetAll called!")

	
	conn, err := dialGrpc()
	if err != nil {
		fmt.Printf("failed to dial: %v\n", err)
		log.Fatalf("failed to dial: %v", err)
		return 
	}
	defer conn.Close()

	client := rundoogrpc.NewProductServiceClient(conn)

	response, err := client.GetProducts(context.Background(), &rundoogrpc.GetProductsRequest{})
	if err != nil {
		fmt.Printf("failed to get products: %v\n", err)
		log.Printf("failed to get products: %v", err)
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
	fmt.Printf("Get products: %v", products)
	return 
}

func (m *RundooModel) Get(id int64) (product data.Product, err error) {
	log.Printf("RundooModel Get %d called!", id)

	conn, err := dialGrpc()
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

func (m *RundooModel) AddProduct(product *data.Product) (err error) {
	log.Println("RundooModel AddProduct called!")

	conn, err := dialGrpc()
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return 
	}
	defer conn.Close()

	client := rundoogrpc.NewProductServiceClient(conn)

	gProduct := rundoogrpc.Product{
		Name: product.Name,
		Category: string(product.Category),
		Sku: string(product.Sku),
	}

	_, err = client.AddProduct(context.Background(), &rundoogrpc.AddProductRequest{Product: &gProduct})
	if err != nil {
		log.Fatalf("failed to add product: %v", err)
		return 
	}

	
	return 
}

func (m *RundooModel) SearchProducts(filters []rundoogrpc.Filter) (products []data.Product, err error) {
	log.Println("RundooModel SearchProducts called!")

	
	conn, err := dialGrpc()
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return 
	}

	client := rundoogrpc.NewProductServiceClient(conn)	
	
	// Create the search request
	filtersp := make([]*rundoogrpc.Filter, len(filters))
	for i, filter := range filters {
		filtersp[i] = &filter
	}

	request := &rundoogrpc.SearchProductsRequest{Filters: filtersp}

	// Invoke the gRPC method
	response, err := client.SearchProducts(context.Background(), request)
	if err != nil {
		log.Fatalf("SearchProducts failed: %v", err)
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

func dialGrpc() (conn *grpc.ClientConn, err error) {
	serviceURL, err := registry.GetProvider(registry.RundooService)
	fmt.Printf("dialGrpc: serviceURL: %s\n", serviceURL)
	if err != nil {
		log.Println("Error getting provider RundooService: ", err)
		fmt.Println("Error getting provider RundooService: ", err)
		return 
	}
	record := strings.Split(serviceURL, ":") // http://localhost:port
	portInt, _ := strconv.Atoi(record[2])
	rpcPort := ":" + strconv.Itoa(portInt+1)
	fmt.Printf("dialGrpc serviceURL %s : %v %s %s \n", serviceURL, record, record[0], record[1])
	fmt.Printf("dialGrpc record[1][2:]+rpcPort %s \n", record[1][2:]+rpcPort)

	//conn, err = grpc.Dial("localhost"+rpcPort, grpc.WithInsecure())
	conn, err = grpc.Dial(record[1][2:]+rpcPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		fmt.Printf("failed to dial: %v", err)
		return 
	}
	return
}



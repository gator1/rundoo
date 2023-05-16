package rundoo

import (
	"log"
	"sync"

	rundoogrpc "app/api/v1"
	"app/internal/data"
)

type CategoryType string

const (
	CategoryWine = CategoryType("Wine")
	CategoryBook = CategoryType("Book")
	CategoryTool = CategoryType("Tool")
)


type ProductService struct {
	// a database dependency would go here but instead we're going to have a static map
	products data.Products
	models *data.Models
    Categories map[CategoryType]bool
    productsMutex sync.Mutex
}

type ServiceInterface interface {
	GetProducts() (data.Products, error)
	SearchProducts(filters []rundoogrpc.Filter) (data.Products, error)
	AddProduct(product data.Product) (bool, error)
}


// NewService instantiates a new Service.
func NewService( models *data.Models) ServiceInterface {
	return &ProductService{
		products: *new(data.Products),
		models: models,
		Categories: map[CategoryType]bool{
			CategoryWine: true,
        	CategoryBook: true,
       		 CategoryTool: true,
		},
		productsMutex: *new(sync.Mutex),
	}
}

var (
    products data.Products
    Categories = map[CategoryType]bool{
        CategoryWine: true,
        CategoryBook: true,
        CategoryTool: true,
    }
    productsMutex sync.Mutex
)

func (s *ProductService) GetProducts() (result data.Products, err error) {
	// instead of querying a database, we just query our static map
	products, err := s.models.Products.GetAll()
	if err != nil {
		log.Printf("ProductService, GetProducts %v", err)
		return
	}
	for _, mproduct := range products {
		result = append(result, data.Product{
			ID: mproduct.ID,
			Name:  mproduct.Name,
			Category:  data.CategoryType(mproduct.Category),
			Sku:  data.SKU(mproduct.Sku),

		})

	}
	return 
}


func (p *ProductService) SearchProducts(filters []rundoogrpc.Filter) (data.Products, error) {
	// Create a new Products slice to store filtered products
	filteredProducts := make(data.Products, 0)

	// Loop over the p.Products slice and apply filters
	for _, product := range p.products {
		if product.MatchFilters(filters) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	// Return the filtered products slice and any error that may have occurred
	return filteredProducts, nil
}



/*
func (s *ProductService) SearchProducts(ctx context.Context, req *rundoogrpc.SearchProductsRequest) (*rundoogrpc.SearchProductsResponse, error) {
	// instead of querying a database, we just query our static map

	var matchedProducts Products
	for _, filter := range req.GetFilters() {
		switch filter.GetField() {
		case "name":
			for _, product := range s.products {
				if product.Name == filter.GetValue() {
					matchedProducts = append(matchedProducts, product)
				}
			}
		case "category":
			for _, product := range s.products {
				if product.Category == CategoryType(filter.GetValue()) {
					matchedProducts = append(matchedProducts, product)
				}
			}
		case "sku":
			sku, err := NewSKU(filter.GetValue())
			if err != nil {
				return nil, err
			}
			product, err := s.products.GetBySKU(sku)
			if err != nil {
				return nil, err
			}
			matchedProducts = append(matchedProducts, *product)
		}
	}

	return &rundoogrpc.SearchProductsResponse{
		Products: matchedProducts.toProto(),
	}, nil
}
*/

func (s *ProductService) AddProduct(product data.Product) (bool,error) {
	// Instead of using a database, we just add the product to our static map.


	return true, nil
}


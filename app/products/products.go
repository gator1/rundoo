package products

import (
	"fmt"
	"regexp"
	"sync"
	rundoogrpc "app/api/v1"
)

type SKU string

func NewSKU(s string) (SKU, error) {
	// Use a regular expression to validate the SKU
	pattern := `^[A-Z]{2}[1-9A-Z]{6,10}$`
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return "", fmt.Errorf("failed to validate SKU: %v", err)
	}
	if !matched {
		return "", fmt.Errorf("invalid SKU format")
	}
	return SKU(s), nil
}

type CategoryType string

const (
	CategoryWine = CategoryType("Wine")
	CategoryBook = CategoryType("Book")
	CategoryTool = CategoryType("Tool")
)

type Product struct {
	Sku      SKU
	Category CategoryType
	Name     string
	
}

type Products []Product

type ProductService struct {
	// a database dependency would go here but instead we're going to have a static map
	products Products
    Categories map[CategoryType]bool
    productsMutex sync.Mutex
}

type ServiceInterface interface {
	GetProducts() (Products, error)
	SearchProducts(filters []rundoogrpc.Filter) (Products, error)
	AddProduct(product Product) (bool, error)
}


// NewService instantiates a new Service.
func NewService( /* a database connection would be injected here */ ) ServiceInterface {
	return &ProductService{
		products: *new(Products),
		Categories: map[CategoryType]bool{
			CategoryWine: true,
        	CategoryBook: true,
       		 CategoryTool: true,
		},
		productsMutex: *new(sync.Mutex),
	}
}

var (
    products Products
    Categories = map[CategoryType]bool{
        CategoryWine: true,
        CategoryBook: true,
        CategoryTool: true,
    }
    productsMutex sync.Mutex
)

func (s *ProductService) GetProducts() (result Products, err error) {
	// instead of querying a database, we just query our static map
	
	return products, nil
}

func (p Products) toProto() []*rundoogrpc.Product {
	protoProducts := make([]*rundoogrpc.Product, len(p))
	for i, product := range p {
		protoProducts[i] = &rundoogrpc.Product{
			Name:     product.Name,
			Category: string(product.Category),
			Sku:      string(product.Sku),
		}
	}
	return protoProducts
}

func (p *ProductService) SearchProducts(filters []rundoogrpc.Filter) (Products, error) {
	// Create a new Products slice to store filtered products
	filteredProducts := make(Products, 0)

	// Loop over the p.Products slice and apply filters
	for _, product := range p.products {
		if product.MatchFilters(filters) {
			filteredProducts = append(filteredProducts, product)
		}
	}

	// Return the filtered products slice and any error that may have occurred
	return filteredProducts, nil
}

func (p *Product) MatchFilters(filters []rundoogrpc.Filter) bool {
	// Apply the filters to the product and determine if it matches
	// Return true if the product matches all filters, false otherwise
	// You would need to implement the logic specific to your application's filter criteria
	return true
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

func (s *ProductService) AddProduct(product Product) (bool,error) {
	// Instead of using a database, we just add the product to our static map.


	return true, nil
}


func (p Products) GetByName(name string) (*Product, error) {
	for i := range p {
		if p[i].Name == name {
			return &p[i], nil
		}
	}

	return nil, fmt.Errorf("Product with Name '%v' not found", name)
}

func (p Products) GetBySKU(sku SKU) (*Product, error) {
	for i := range p {
		if p[i].Sku == sku {
			return &p[i], nil
		}
	}

	return nil, fmt.Errorf("Product with Sku '%v' not found", sku)
}
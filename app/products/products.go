package products

import (
	"fmt"
	"regexp"
	"sync"
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
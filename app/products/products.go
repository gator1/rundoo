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
	Name     string
	Category CategoryType
	Sku      SKU
}

type Products []Product

var (
    products Products
    categories = map[CategoryType]bool{
        CategoryWine: true,
        CategoryBook: true,
        CategoryTool: true,
    }
    productsMutex sync.Mutex
)

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
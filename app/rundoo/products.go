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
	products data.Products
	models *data.Models
    Categories map[CategoryType]bool
    productsMutex sync.Mutex
}

type ServiceInterface interface {
	GetProducts() (data.Products, error)
	GetProduct(id int64) (data.Product, error)
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

func (s *ProductService) GetProduct(id int64) (result data.Product, err error) {
	product, err := s.models.Products.Get(id)
	if err != nil {
		log.Printf("ProductService, GetProduct %d %v", id, err)
		return
	}
	result.ID = product.ID
	result.Name = product.Name
	result.Category = data.CategoryType(product.Category)
	result.Sku =  data.SKU(product.Sku)

	return 
}



func (s *ProductService) SearchProducts(filters []rundoogrpc.Filter) (products data.Products, err error) {
	
	rundooproducts, err := s.models.Products.SearchProducts(filters)
	if err != nil {
		log.Printf("ProductService, SearchProducts %v", err)
		return
	}
	
	// Create a new Products slice to store filtered products
	filteredProducts := make(data.Products, 0)

	// Loop over the p.Products slice and apply filters
	for _, product := range rundooproducts {
		filteredProducts = append(filteredProducts, *product)
		
	}

	// Return the filtered products slice and any error that may have occurred
	return filteredProducts, nil
	
}

func (s *ProductService) AddProduct(product data.Product) (result bool, err error) {
    err = s.models.Products.Insert(&product)
	if err != nil {
		log.Printf("ProductService, AddProduct %v", err)
		return
	}

	return true, nil
}


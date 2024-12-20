package data

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/lib/pq"

	rundoogrpc "app/api/v1"
)

type SKU string

type CategoryType string

const (
	CategoryWine = CategoryType("Wine")
	CategoryBook = CategoryType("Book")
	CategoryTool = CategoryType("Tool")
)

type Product struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string
	Category  CategoryType
	Sku       SKU
	Version   int32 `json:"-"`
}

type Products []Product

type ProductModel struct {
	DB *sql.DB
}

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

func (b ProductModel) Insert(product *Product) error {
	query := `
		INSERT INTO products (name, category, sku)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version`

	args := []interface{}{product.Name, product.Category, product.Sku}
	// return the auto generated system values to Go object
	return b.DB.QueryRow(query, args...).Scan(&product.ID, &product.CreatedAt, &product.Version)
}

func (b ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	query := `
		SELECT id, created_at, name, category, sku, version
		FROM products
		WHERE id = $1`

	var product Product

	err := b.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.Name,
		&product.Category,
		&product.Sku,
		&product.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (b ProductModel) Update(product *Product) error {
	query := `
		UPDATE products
		SET name = $1, category = $2, sku = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []interface{}{product.Name, product.Category, product.Sku, product.ID, product.Version}
	return b.DB.QueryRow(query, args...).Scan(&product.Version)
}

func (b ProductModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	query := `
		DELETE FROM products
		WHERE id = $1`

	results, err := b.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (b ProductModel) GetAll() ([]*Product, error) {
	query := `
	  SELECT * 
	  FROM products
	  ORDER BY id`

	rows, err := b.DB.Query(query)
	if err != nil {
		fmt.Printf("ProductModel GetAll err %v\n", err)
		return nil, err
	}

	defer rows.Close()

	products := []*Product{}

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.ID,
			&product.CreatedAt,
			&product.Name,
			&product.Category,
			&product.Sku,
			&product.Version,
		)
		if err != nil {
			fmt.Printf("ProductModel GetAll row scan err %v\n", err)
			return nil, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("ProductModel GetAll row err %v\n", err)
		return nil, err
	}

	fmt.Printf("ProductModel GetAll %v\n", products)
	return products, nil
}

func (b ProductModel) SearchProducts(filters []rundoogrpc.Filter) ([]*Product, error) {
	// Perform the search query
	query := `
	SELECT * FROM Products
	WHERE name ILIKE $1 OR category ILIKE $1 OR sku ILIKE $1
`

	// Define the search parameters
	searchString := "%search_string%"
	includeName := false
	includeCategory := false
	includeSKU := false

	for _, filter := range filters {
		if filter.Field == "Any" {
			includeName = true
			includeCategory = true
			includeSKU = true
			searchString = filter.Value
			break
		}
		if filter.Field == "Name" {
			includeName = true
			searchString = filter.Value
		}
		if filter.Field == "Category" {
			includeCategory = true
			searchString = filter.Value
		}
		if filter.Field == "Sku" {
			includeSKU = true
			searchString = filter.Value
		}

	}

	// Construct the search query dynamically
	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("SELECT * FROM Products WHERE ")

	params := make([]interface{}, 0)
	paramIndex := 1

	if includeName {
		queryBuilder.WriteString(fmt.Sprintf("name ILIKE $%d", paramIndex))
		params = append(params, searchString)
		paramIndex++
	}

	if includeCategory {
		if includeName {
			queryBuilder.WriteString(" OR ")
		}
		queryBuilder.WriteString(fmt.Sprintf("category ILIKE $%d", paramIndex))
		params = append(params, searchString)
		paramIndex++
	}

	if includeSKU {
		if includeName || includeCategory {
			queryBuilder.WriteString(" OR ")
		}
		queryBuilder.WriteString(fmt.Sprintf("sku ILIKE $%d", paramIndex))
		params = append(params, searchString)
		paramIndex++
	}

	query = queryBuilder.String()

	// Perform the search query
	rows, err := b.DB.Query(query, params...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterate over the query results

	products := []*Product{}

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.ID,
			&product.CreatedAt,
			&product.Name,
			&product.Category,
			&product.Sku,
			&product.Version,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p Product) MatchFilters(filters []rundoogrpc.Filter) bool {
	// Apply the filters to the product and determine if it matches
	// Return true if the product matches all filters, false otherwise
	// You would need to implement the logic specific to your application's filter criteria
	return true
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

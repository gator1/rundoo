package rundooportal

import (
	"app/products"
	"app/registry"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func RegisterHandlers() {
	http.Handle("/", http.RedirectHandler("/products", http.StatusPermanentRedirect))

	h := new(productsHandler)
	http.Handle("/products", h)
	http.Handle("/products/", h)
}

type productsHandler struct{}

func (sh productsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("rundooportal productsHandler Request received", r.URL.Path)
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2: // /products
		sh.renderProducts(w, r)
	case 3: // /products/{:sku}
		sku := products.SKU(pathSegments[2])
		
		sh.renderProduct(w, r, sku)

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (productsHandler) renderProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving products: ", err)
		}
	}()

	serviceURL, err := registry.GetProvider(registry.ProductService)
	if err != nil {
		log.Println("Error getting provider ProductService: ", err)
		return
	}

	res, err := http.Get(serviceURL + "/products")
	if err != nil {
		log.Println("Error http get peoducts: ", err)
		return
	}

	var s products.Products
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Println("Error json decodes peoducts: ", err)
		return
	}

	rootTemplate.Lookup("products.gohtml").Execute(w, s)
}

func (productsHandler) renderProduct(w http.ResponseWriter, r *http.Request, sku products.SKU) {

	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving products: ", err)
			return
		}
	}()

	serviceURL, err := registry.GetProvider(registry.ProductService)
	if err != nil {
		return
	}

	res, err := http.Get(fmt.Sprintf("%v/products/%v", serviceURL, string(sku)))
	if err != nil {
		log.Println("Error request product : ", string(sku))
		return
	}

	var s products.Product
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Println("Error decodes product : ", string(sku))
		return
	}

	rootTemplate.Lookup("addproduct.gohtml").Execute(w, s)
}

func (productsHandler) renderGrades(w http.ResponseWriter, r *http.Request, id int) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		w.Header().Add("location", fmt.Sprintf("/products/%v", id))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()
	name := r.FormValue("Name")
	categoryType := r.FormValue("Category")
	sku := r.FormValue("Sku")
	g := products.Product{
		Name: name,
		Category:  products.CategoryType(categoryType),
		Sku: products.SKU(sku),
	}
	data, err := json.Marshal(g)
	if err != nil {
		log.Println("Failed to convert product to JSON: ", g, err)
	}

	serviceURL, err := registry.GetProvider(registry.ProductService)
	if err != nil {
		log.Println("Failed to retrieve instance of Product Service", err)
		return
	}
	res, err := http.Post(fmt.Sprintf("%v/products/%v", serviceURL, name), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to save product to Products Service", err)
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.Println("Failed to save product to Product Service. Status: ", res.StatusCode)
		return
	}
}

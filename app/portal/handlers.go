package rundooportal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	rundoogrpc "app/api/v1"
	"app/registry"
	"app/internal/data"
)

func HttpHandler() {
	http.Handle("/", http.RedirectHandler("/products", http.StatusPermanentRedirect))

	h := new(RundooHandler)
	http.Handle("/products", h)
	http.Handle("/products/", h)

	
}

type RundooHandler struct{}

func (sh RundooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("rundooportal productsHandler Request received", r.URL.Path)
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2: // /products
		sh.renderProductsGrpc(w, r)
	case 3: // /products/{:sku}
		sku := data.SKU(pathSegments[2])
		if sku == "AddProduct" {
			sh.renderAddProduct(w, r)

		} else if sku == "AddedProduct" {
			sh.postAddProduct(w, r)

		} else if sku == "SearchProduct" {

		} else {
			sh.renderProduct(w, r, sku)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (RundooHandler) renderProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving products: ", err)
		}
	}()

	serviceURL, err := registry.GetProvider(registry.RundooService)
	if err != nil {
		log.Println("Error getting provider ProductService: ", err)
		return
	}

	res, err := http.Get(serviceURL + "/products")
	if err != nil {
		log.Println("Error http get peoducts: ", err)
		return
	}

	var s data.Products
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Println("Error json decodes peoducts: ", err)
		return
	}

	err = rootTemplate.ExecuteTemplate(w, "base", s)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	// rootTemplate.Lookup("home.html").Execute(w, s)
}

func (RundooHandler) renderProductsGrpc(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving products: ", err)
		}
	}()

	log.Println("renderProductsGrpc called!")

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
	}
	defer conn.Close()

	client := rundoogrpc.NewProductServiceClient(conn)

	response, err := client.GetProducts(context.Background(), &rundoogrpc.GetProductsRequest{})
	if err != nil {
		log.Fatalf("failed to get products: %v", err)
	}

	for _, product := range response.GetProducts() {
		log.Printf("Product: %v\n", product)
	}
	err = rootTemplate.ExecuteTemplate(w, "base", response.Products)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	// rootTemplate.Lookup("products.gohtml").Execute(w, response.Products)
}



func (RundooHandler) renderProduct(w http.ResponseWriter, r *http.Request, sku data.SKU) {

	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error retrieving products: ", err)
			return
		}
	}()

	serviceURL, err := registry.GetProvider(registry.RundooService)
	if err != nil {
		log.Println("Error redner product GetProvider : ", string(sku))

		return
	}

	res, err := http.Get(fmt.Sprintf("%v/products/%v", serviceURL, string(sku)))
	if err != nil {
		log.Println("Error render product http.Get: ", string(sku))
		return
	}

	var s data.Product
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Println("Error decodes product : ", string(sku))
		return
	}

	rootTemplate.Lookup("productdetails.gohtml").Execute(w, s)
}

func (RundooHandler) renderAddProduct(w http.ResponseWriter, r *http.Request) {

	rootTemplate.Lookup("addproduct.gohtml").Execute(w, nil)

	defer func() {
		w.Header().Add("location", "/products")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()
}

func (RundooHandler) postAddProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		w.Header().Add("location", "/products")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}()

	serviceURL, err := registry.GetProvider(registry.RundooService)
	if err != nil {
		log.Println("Error redner product GetProvider : ", err)

		return
	}

	name := r.FormValue("Name")
	categoryType := data.CategoryType(r.FormValue("Category"))
	sku, err := data.NewSKU(r.FormValue("Sku"))
	if err != nil {
		log.Println("Wrong format of SKU: ", err)
		return
	}

	p := data.Product{
		Name:     name,
		Category: categoryType,
		Sku:      sku,
	}

	data, err := json.Marshal(p)
	if err != nil {
		log.Println("Failed to convert product to JSON: ", p, err)
	}

	res, err := http.Post(fmt.Sprintf("%v/products/AddProduct", serviceURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to save product to Product Service", err)
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.Println("Failed to save product to Product Service. Status: ", res.StatusCode)
		return
	}
}

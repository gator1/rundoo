package rundoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HttpHandler() {
	handler := new(ProductsHandler)
	http.Handle("/products", handler)
	http.Handle("/products/", handler)
}

type ProductsHandler struct{}

func (sh ProductsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2: // /products
		sh.getAll(w, r)
	
	case 3: // /products/{:sku}
		sku := pathSegments[2]
		if sku == "AddProduct" {
			sh.addProduct(w, r)
		} else {
			log.Println("demo for central log service: we don't have a way to go to an indivisual product")
			sh.getOne(w, r, SKU(sku))
		}
	
		
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh ProductsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	data, err := sh.toJSON(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (sh ProductsHandler) getOne(w http.ResponseWriter, r *http.Request, sku SKU) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	product, err := products.GetBySKU(sku)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}
	}

	data, err := sh.toJSON(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("Failed to serialize products: %q", err))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (ProductsHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize products: %q", err)
	}
	return b.Bytes(), nil
}



func (sh ProductsHandler) addProduct(w http.ResponseWriter, r *http.Request) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	var g Product
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	products = append(products, g)

	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}


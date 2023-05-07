package products

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func RegisterHandlers() {
	handler := new(productsHandler)
	http.Handle("/products", handler)
	http.Handle("/products/", handler)
}

type productsHandler struct{}

func (sh productsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2: // /product
		sh.getAll(w, r)
	
	case 3: // /products/{:sku}
		log.Println("demo for central log service: we don't have a way to go to an indivisual product")

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh productsHandler) getAll(w http.ResponseWriter, r *http.Request) {
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

func (sh productsHandler) getOne(w http.ResponseWriter, r *http.Request, sku SKU) {
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

func (productsHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize products: %q", err)
	}
	return b.Bytes(), nil
}

/*

func (sh productsHandler) addGrade(w http.ResponseWriter, r *http.Request, sku SKU) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	student, err := products.GetBySKU(sku)
	if err != nil {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}

	var g Product
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	student.Grades = append(student.Grades, g)

	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}
*/

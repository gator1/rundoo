package rundoo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"app/internal/data"
)

func HttpHandler() {
	handler := new(ProductsHandler)
	http.Handle("/products", handler)
	http.Handle("/products/", handler)
}

type ProductsHandler struct{}

type Envelope map[string]any

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
			if sku == "" {
				sh.getAll(w, r)
			} else {

			
			log.Printf("demo for central log service: we don't have a way to go to an indivisual product %s", sku)
			sh.getOne(w, r, data.SKU(sku))
		}
		}
	
		
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh ProductsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	data, err := sh.ToJSON(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (sh ProductsHandler) getOne(w http.ResponseWriter, r *http.Request, sku data.SKU) {
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

	data, err := sh.ToJSON(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(fmt.Errorf("Failed to serialize products: %q", err))
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}

func (ProductsHandler) ToJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize products: %q", err)
	}
	return b.Bytes(), nil
}


// Credit: Alex Edwards, Let's Go Further
func (ProductsHandler) WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}


func (ProductsHandler) ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		// Custom Error Handling: Alex Edwards, Let's Go Further Chapter 4
		return err
	}

	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}



func (sh ProductsHandler) addProduct(w http.ResponseWriter, r *http.Request) {
	productsMutex.Lock()
	defer productsMutex.Unlock()

	var g data.Product
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	products = append(products, g)

	w.WriteHeader(http.StatusCreated)
	data, err := sh.ToJSON(g)
	if err != nil {
		log.Println(err)
	}
	w.Header().Add("content-type", "application/json")
	w.Write(data)
}


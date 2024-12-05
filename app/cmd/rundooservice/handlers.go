package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"app/internal/data"
	"app/rundoo"
)


func (app *application) getCreateProductsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		products, err := app.models.Products.GetAll()
		if err != nil {
			fmt.Printf("failed to get products: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}


		data, err := app.handler.ToJSON(products)
		if err != nil {
			fmt.Printf("failed ToJSON: %v\n", err)
			
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		w.Header().Add("content-type", "application/json")
		w.Write(data)
	}

	if r.Method == http.MethodPost {
		var input struct {
			Name     string   `json:"name"`
			Category string   `json:"category"`
			Sku      string   `json:"sku"`
		}

		err := app.handler.ReadJSON(w, r, &input)
		if err != nil {
			fmt.Printf("failed ReadJSON: %v\n", err)
			
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		product := &data.Product{
			Name:     input.Name,
			Category: data.CategoryType(input.Category),
			Sku:      data.SKU(input.Sku), 
		}

		err = app.models.Products.Insert(product)
		if err != nil {
			fmt.Printf("failed Products Insert: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("products/%d", product.ID))

		// Write the JSON response with a 201 Created status code and the Location header set.
		err = app.handler.WriteJSON(w, http.StatusCreated, rundoo.Envelope{"product": product}, headers)
		if err != nil {
			fmt.Printf("failed Products WriteJSON: %v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func (app *application) getUpdateDeleteProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getProduct(w, r)
	case http.MethodPut:
		app.updateProduct(w, r)
	case http.MethodDelete:
		app.deleteProduct(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/products/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	product, err := app.models.Products.Get(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if err := app.handler.WriteJSON(w, http.StatusOK, rundoo.Envelope{"product": product}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/books/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	product, err := app.models.Products.Get(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Name     *string     `json:"name"`
		Category *string     `json:"category"`
		Sku     *string      `json:"sku"`
	}

	err = app.handler.ReadJSON(w, r, &input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}

	if input.Category != nil {
		product.Category = data.CategoryType(*input.Category)
	}

	if input.Sku != nil {
		product.Sku = data.SKU(*input.Sku)
	}

	err = app.models.Products.Update(product)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := app.handler.WriteJSON(w, http.StatusOK, rundoo.Envelope{"product": product}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/products/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = app.models.Products.Delete(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = app.handler.WriteJSON(w, http.StatusOK, rundoo.Envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

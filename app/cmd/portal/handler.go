package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	rundoogrpc "app/api/v1"
	"app/internal/data"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	products, err := app.productlist.GetAll()
	
	if err != nil {
		fmt.Println("home err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		fmt.Println("home template template", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", products)
	if err != nil {
		log.Print(err.Error())
		fmt.Println("home template ExecuteTemplate", err)
		http.Error(w, "Internal server error", 500)
		return
	}


}

func (app *application) productView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		fmt.Println("productView, get id not found", err)
		http.NotFound(w, r)
		return
	}

	product, err := app.productlist.Get(int64(id))
	if err != nil {
		fmt.Println("productView, get product not found", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/view.html",
	}

	// Used to convert comma-separated genres to a slice within the template.
	funcs := template.FuncMap{"join": strings.Join}

	ts, err := template.New("showProduct").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		fmt.Println("productView, show product", err)
		
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", product)
	if err != nil {
		fmt.Println("productView,ExecuteTemplate", err)
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) productCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.productCreateForm(w, r)
	case http.MethodPost:
		app.productCreateProcess(w, r)
	default:
		fmt.Println("productView,productCreate, not allowed")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) productsSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.productsSearchForm(w, r)
	case http.MethodPost:
		app.productsSearchProcess(w, r)
	default:
		fmt.Println("productView,productsSearch, not allowed")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}



func (app *application) productCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/create.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("productCreateForm, ParseFiles err", err)
		
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println("productCreateForm, ExecuteTemplate err", err)
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) productCreateProcess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("productCreateProcess, ParseForm err", err)
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.PostForm.Get("name")

	category := r.PostForm.Get("category") 
	sku := r.PostForm.Get("sku") 
	
	product := data.Product {
		Name:     name,
		Category:     data.CategoryType(category),
		Sku: data.SKU(sku),
	}

	err = app.productlist.AddProduct(&product)
	if err != nil {
		fmt.Println("productCreateProcess, AddProduct err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	
}


func (app *application) productsSearchForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/search.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("productsSearchForm, ParseFiles err", err)
		
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println("productsSearchForm, ParExecuteTemplateseFiles err", err)
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) productsSearchProcess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("productsSearchProcess, ParseForm err", err)
		log.Println(err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	searchQuery := r.PostForm.Get("q")
	filterType := r.PostForm.Get("Type")
	
	filters := []rundoogrpc.Filter{{Field: filterType, Value: searchQuery}}
	
	products, err := app.productlist.SearchProducts(filters)
	if err != nil {
		fmt.Println("productsSearchProcess, SearchProducts err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("productsSearchProcess, ParseFiles err", err)
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", products)
	if err != nil {
		fmt.Println("productsSearchProcess, ExecuteTemplate err", err)
		log.Print(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	
}

package rundooportal

import (
	"fmt"
	"html/template"
	"os"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	
	cwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        return err
    }
    fmt.Println("Current working directory:", cwd)

	rootTemplate, err = template.ParseFiles(
		"portal/products.gohtml",
		"portal/addproduct.gohtml",
		"portal/productdetails.gohtml")

	if err != nil {
		return err
	}

	return nil
}
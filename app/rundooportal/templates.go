package rundooportal

import (
	"html/template"
)

var rootTemplate *template.Template

func ImportTemplates() error {
	var err error
	rootTemplate, err = template.ParseFiles(
		"rundooportal/products.gohtml",
		"rundooportal/addproduct.gohtml")

	if err != nil {
		return err
	}

	return nil
}

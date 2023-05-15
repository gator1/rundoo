package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Product struct {
	ID        int64    `json:"id"`
	Name     string   `json:"name"`
	Category int      `json:"category"`
	Sku      int      `json:"sku"`
}

type ProductResponse struct {
	Book *Product `json:"product"`
}

type ProductsResponse struct {
	products *[]Product `json:"products"`
}

type RundooModel struct {
	Endpoint string
}

func (m *RundooModel) GetAll() (*[]Product, error) {
	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var productsResp ProductsResponse
	err = json.Unmarshal(data, &productsResp)
	if err != nil {
		return nil, err
	}

	return productsResp.products, nil
}

func (m *RundooModel) Get(id int64) (*Product, error) {
	url := fmt.Sprintf("%s/%d", m.Endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var productResp ProductResponse
	err = json.Unmarshal(data, &productResp)
	if err != nil {
		return nil, err
	}

	return productResp.Book, nil
}

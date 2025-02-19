package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"user-favorites-service/internals/models"
)

type ProductClient struct {
	URL string
}

func NewProductClient(host string) *ProductClient {
	URL := host
	return &ProductClient{URL: URL}
}

func (productClient *ProductClient) GetProductById(productID int) (models.ProductResponse, error) {
	productsData := models.ProductResponse{}
	url := fmt.Sprintf("%s/products/%d", productClient.URL, productID)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("err from get isteği", err)
		return models.ProductResponse{}, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err from body okuma isteği", err)
		return models.ProductResponse{}, err
	}

	err = json.Unmarshal(bytes, &productsData)
	if err != nil {
		fmt.Println("err unmarshall için", err)
		return models.ProductResponse{}, err
	}
	defer resp.Body.Close()

	return productsData, nil
}

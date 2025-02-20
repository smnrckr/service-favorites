package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"user-favorites-service/internals/models"

	"github.com/pkg/errors"
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
	var resp *http.Response
	maxAttemped := 5

	for i := 1; i <= maxAttemped; i++ {
		response, err := http.Get(url)
		if err != nil {
			return models.ProductResponse{}, err
		}
		if response.StatusCode < 500 {
			resp = response
			break
		}

		time.Sleep(2000)

	}

	if resp == nil {
		return models.ProductResponse{}, errors.New("Max attemped reached")
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ProductResponse{}, err
	}

	err = json.Unmarshal(bytes, &productsData)
	if err != nil {
		return models.ProductResponse{}, err
	}

	defer resp.Body.Close()

	return productsData, nil
}

package models

type ProductResponse struct {
	Id           int    `json:"id"`
	ProductName  string `json:"product_name"`
	ProductCode  string `json:"product_code"`
	ProductPrice string `json:"product_price"`
}

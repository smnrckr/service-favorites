package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Favorite struct {
	Id        int `json:"id" gorm:"column:id"`
	UserID    int `json:"user_id" gorm:"column:user_id"`
	ProductID int `json:"product_id" gorm:"column:product_id"`
	ListID    int `json:"list_id" gorm:"column:list_id"`
}

type FavoriteCreateRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	ListID    int `json:"list_id"`
}

func (f FavoriteCreateRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.UserID, validation.Required.Error("userId bulunmalı")),
		validation.Field(&f.ProductID, validation.Required.Error("productId bulunmalı")),
		validation.Field(&f.ListID, validation.Required.Error("listId bulunmalı")),
	)
}

func (Favorite) TableName() string {
	return "favorite"
}

type FavoritesResponse struct {
	Product []ProductResponse `json:"products"`
}

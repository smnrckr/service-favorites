package models

import validation "github.com/go-ozzo/ozzo-validation"

type FavoriteList struct {
	Id     int    `json:"id" gorm:"column:id"`
	UserID int    `json:"user_id" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:name"`
}

type FavoriteListCreateRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func (f FavoriteListCreateRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Length(1, 50).Error("name uzunluğu 1 - 50 arasında olmalı"), validation.Required.Error("name alanı zorunlu")),
		validation.Field(&f.UserID, validation.Required),
	)
}

type FavoriteListUpdateRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func (f FavoriteListUpdateRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.Length(1, 50).Error("name uzunluğu 1 - 50 arasında olmalı"), validation.Required.Error("name alanı zorunlu")),
	)
}

func (FavoriteList) TableName() string {
	return "favorite_list"
}

type FavoriteListResponse struct {
	FavoriteList []FavoriteList `json:"favorite_list"`
}

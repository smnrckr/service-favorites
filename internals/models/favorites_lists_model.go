package models

type FavoriteList struct {
	Id     int    `json:"id" gorm:"column:id"`
	UserID int    `json:"user_id" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:name"`
}

type FavoriteListCreateRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type FavoriteListUpdateRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

func (FavoriteList) TableName() string {
	return "favorite_list"
}

type FavoriteListResponse struct {
	FavoriteList []FavoriteList `json:"favorite_list"`
}

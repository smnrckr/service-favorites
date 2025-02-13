package models

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

func (Favorite) TableName() string {
	return "favorite"
}

type FavoritesResponse struct {
	Favorite []Favorite `json:"favorites"`
}

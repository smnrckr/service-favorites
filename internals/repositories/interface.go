package repositories

import "gorm.io/gorm"

type PostgreDB interface {
	GetConnection() *gorm.DB
	Close()
}

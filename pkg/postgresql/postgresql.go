package postgresql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	Dbname     string
	Dbuser     string
	Dbpassword string
	Host       string
	Port       string
}

type DB struct {
	DB *gorm.DB
}

func NewDB(config DbConfig) *DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.Host, config.Dbuser, config.Dbpassword, config.Dbname, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	fmt.Println("connection successful")
	return &DB{db}
}

func (db *DB) GetConnection() *gorm.DB {
	return db.DB

}

func (db *DB) Close() {
	sqlDB, err := db.DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()

}

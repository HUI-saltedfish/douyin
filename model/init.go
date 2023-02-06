package model

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error
	dsn := "root:123456@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate
	DB.AutoMigrate(&User{})

	return DB, nil
}

func GetDB() (*gorm.DB, error) {
	var once sync.Once
	once.Do(func() {
		_, err := InitDB()
		if err != nil {
			panic(err)
		}
	})
	return DB, nil
}

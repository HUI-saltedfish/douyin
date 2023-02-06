package model

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error
	dsn := "douyin:123456@tcp(114.116.80.86:3306)/douyin?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// delete table
	// DB.Migrator().DropTable(&User{}, &Video{}, &Comment{})

	// AutoMigrate
	DB.AutoMigrate(&User{}, &Video{}, &Comment{})

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

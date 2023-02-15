package model

import (
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func InitDB() (*gorm.DB, error) {
	var err error
	dsn := "root:123456@tcp(www.huilearn.work:3340)/douyin?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// AutoMigrate if nessacary
	err = DB.Find(&User{}, &Video{}, &Comment{}).Error
	if err != nil {
		err = DB.AutoMigrate(&User{}, &Video{}, &Comment{})
		if err != nil {
			panic(err)
		}
	}

	return DB, nil
}

func GetDB() (*gorm.DB, error) {
	once.Do(func() {
		_, err := InitDB()
		if err != nil {
			panic(err)
		}
	})
	return DB, nil
}

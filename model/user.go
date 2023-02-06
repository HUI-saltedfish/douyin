package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string `gorm:"type:varchar(20);not null"`
	Password       string `gorm:"type:varchar(20);not null"`
	Follow_count   int    `gorm:"type:int;not null;default:0"`
	Follower_count int    `gorm:"type:int;not null;default:0"`
	Is_follow      bool   `gorm:"type:bool;not null;default:false"`
}

func GetUserByName(name string) (*User, error) {
	var user *User
	db, _ := GetDB()
	err := db.Where("name = ?", name).First(&user).Error
	return user, err
}

func GetUserByNameAndPassword(name string, password string) (*User, error) {
	var user *User
	db, _ := GetDB()
	err := db.Where("name = ? AND password = ?", name, password).First(&user).Error
	return user, err
}

func CreateUser(user *User) error {
	db, _ := GetDB()
	err := db.Create(user).Error
	return err
}

func GetUserById(id int) (*User, error) {
	var user *User
	db, _ := GetDB()
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

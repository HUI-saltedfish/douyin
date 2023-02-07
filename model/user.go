package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name            string    `gorm:"type:varchar(20);not null" json:"name,omitempty"`
	Password        string    `gorm:"type:varchar(20);not null"`
	Follow_count    int       `gorm:"type:int;not null;default:0" json:"follow_count,omitempty"`
	Follower_count  int       `gorm:"type:int;not null;default:0" json:"follower_count,omitempty"`
	Is_follow       bool      `gorm:"type:bool;not null;default:false" json:"is_follow,omitempty"`
	Fllow_Users     []User    `gorm:"many2many:follow_follows;"`
	Has_Videos      []Video   `gorm:"foreignkey:AuthorID"`
	Favorite_Videos []Video   `gorm:"many2many:favorite_videos;"`
	Comments        []Comment `gorm:"foreignkey:UserID"`
	Friends         []User    `gorm:"many2many:user_friends;"`
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

func GetVideosByUserId(id int) ([]Video, error) {
	var videos []Video
	db, _ := GetDB()
	err := db.Where("author_id = ?", id).Find(&videos).Error
	return videos, err
}
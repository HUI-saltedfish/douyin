package model

import (
	"encoding/json"
	"foundation/redisService"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name            string    `gorm:"type:varchar(20);not null; unique"`
	Password        string    `gorm:"type:varchar(20);not null"`
	Follow_count    int       `gorm:"type:int;not null;default:0"`
	Follower_count  int       `gorm:"type:int;not null;default:0"`
	Is_follow       bool      `gorm:"-"`
	Fllow_Users     []User    `gorm:"many2many:follow_follows;"`
	Has_Videos      []Video   `gorm:"foreignkey:AuthorID"`
	Favorite_Videos []Video   `gorm:"many2many:favorite_videos;"`
	Comments        []Comment `gorm:"foreignkey:UserID"`
	Friends         []User    `gorm:"many2many:user_friends;"`
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func GetUserByName(name string) (*User, error) {
	var user = new(User)
	var err error
	// redis first
	err = redisService.RedisClient.HGet(redisService.Ctx, "user", name).Scan(user)
	if err == nil {
		return user, nil
	}

	// mysql second
	db, _ := GetDB()
	err = db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "user", name, user).Err()
	if err != nil {
		return nil, err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "user", 5*time.Minute).Err()
	return user, err
}

func CreateUser(user *User) error {
	db, _ := GetDB()
	err := db.Create(user).Error
	if err != nil {
		return err
	}
	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "user", user.Name, user).Err()
	if err != nil {
		return err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "user", 5*time.Minute).Err()
	return err
}

func GetUserById(id int) (*User, error) {
	var user = new(User)
	var err error

	// redis first
	err = redisService.RedisClient.HGet(redisService.Ctx, "user", strconv.Itoa(id)).Scan(user)
	if err == nil {
		return user, nil
	}

	// mysql second
	db, _ := GetDB()
	err = db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "user", strconv.Itoa(id), user).Err()
	if err != nil {
		return nil, err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "user", 5*time.Minute).Err()
	return user, err
}

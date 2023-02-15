package model

import (
	"encoding/json"
	"interaction/redisService"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Author        User      `gorm:"-"`
	AuthorID      uint      `gorm:"index"`
	PlayUrl       string    `gorm:"type:varchar(100);not null"`
	CoverUrl      string    `gorm:"type:varchar(100);not null"`
	FavoriteCount int       `gorm:"default:0"`
	CommentCount  int       `gorm:"default:0"`
	Is_favorite   bool      `gorm:"-"`
	Title         string    `gorm:"type:varchar(255);not null"`
	Favorite_User []User    `gorm:"many2many:favorite_videos;"`
	Comments      []Comment `gorm:"foreignKey:VideoID"`
}

func (v Video) MarshalBinary() ([]byte, error) {
	return json.Marshal(v)
}

func (v *Video) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, v)
}

func CreateVideo(video *Video) error {
	var err error
	db, _ := GetDB()
	err = db.Create(video).Error
	if err != nil {
		return err
	}
	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "video", strconv.Itoa(int(video.ID)), video).Err()
	if err != nil {
		return err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "video", 5*time.Minute).Err()
	return err

}

func GetVideoOrderByTime(time time.Time) ([]Video, error) {
	db, _ := GetDB()
	var videos []Video
	// get all videos
	err := db.Where("created_at < ?", time).Order("created_at desc").Limit(30).Find(&videos).Error
	return videos, err
}

func GetVideoById(id uint) (Video, error) {
	var video Video
	// redis first
	err := redisService.RedisClient.HGet(redisService.Ctx, "video", strconv.Itoa(int(id))).Scan(&video)
	if err == nil {
		return video, nil
	}

	// if not in redis, get from mysql
	db, _ := GetDB()
	// get video
	err = db.Where("id = ?", id).First(&video).Error
	if err != nil {
		return video, err
	}
	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "video", strconv.Itoa(int(video.ID)), video).Err()
	if err != nil {
		return video, err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "video", 5*time.Minute).Err()
	return video, err
}

func GetVideosByUser(user *User) ([]Video, error) {
	var videos []Video
	db, _ := GetDB()
	err := db.Model(user).Association("Has_Videos").Find(&videos)
	return videos, err
}

func AddFavoriteVideo(user *User, video *Video) error {
	var err error
	db, _ := GetDB()
	err = db.Model(user).Association("Favorite_Videos").Append(video)
	if err != nil {
		return err
	}
	err = UpdateVideoFavoriteCount(video)
	return err
}

func UnFavoriteVideo(user *User, video *Video) error {
	var err error
	db, _ := GetDB()
	err = db.Model(user).Association("Favorite_Videos").Delete(video)
	if err != nil {
		return err
	}
	err = UpdateVideoFavoriteCount(video)
	return err
}

func GetUserFavoriteVideos(user *User) ([]Video, error) {
	var videos []Video
	db, _ := GetDB()
	err := db.Model(user).Association("Favorite_Videos").Find(&videos)
	return videos, err
}

func UpdateVideoFavoriteCount(video *Video) error {
	db, _ := GetDB()
	num_favorite := db.Model(video).Association("Favorite_User").Count()
	return db.Model(video).Update("favorite_count", num_favorite).Error
}

func IsFavoriteVideo(user *User, video *Video) bool {
	db, _ := GetDB()
	var u User
	db.Model(video).Association("Favorite_User").Find(&u, user.ID)
	return u.ID != 0
}

func UpdateVideoCommentCount(video *Video) error {
	db, _ := GetDB()
	num_comment := db.Model(video).Association("Comments").Count()
	return db.Model(video).Update("comment_count", num_comment).Error
}

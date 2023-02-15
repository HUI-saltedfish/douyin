package model

import (
	"encoding/json"
	"foundation/redisService"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	User    User   `gorm:"-"`
	UserID  uint   `gorm:"not null"`
	V       Video  `gorm:"-"`
	VideoID uint   `gorm:"not null"`
}

func (c Comment) MarshalBinary() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Comment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func CreateComment(comment *Comment) error {
	db, _ := GetDB()
	err := db.Create(comment).Error
	if err != nil {
		return err
	}
	video, err := GetVideoById(uint(comment.VideoID))
	if err != nil {
		return err
	}
	err = UpdateVideoCommentCount(&video)
	if err != nil {
		return err
	}

	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "comment", strconv.Itoa(int(comment.ID)), comment).Err()
	if err != nil {
		return err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "comment", 5*time.Minute).Err()
	return err
}

func GetCommentById(id int64) (*Comment, error) {
	var comment Comment
	var err error
	// get from redis
	err = redisService.RedisClient.HGet(redisService.Ctx, "comment", strconv.Itoa(int(id))).Scan(&comment)
	if err == nil {
		return &comment, nil
	}

	// get from mysql
	db, _ := GetDB()
	err = db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}

	// save to redis and expire in 5 minutes
	err = redisService.RedisClient.HSet(redisService.Ctx, "comment", strconv.Itoa(int(id)), comment).Err()
	if err != nil {
		return nil, err
	}
	err = redisService.RedisClient.Expire(redisService.Ctx, "comment", 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func DeleteComment(comment *Comment) error {
	db, _ := GetDB()
	err := db.Delete(comment).Error
	if err != nil {
		return err
	}
	video, err := GetVideoById(uint(comment.VideoID))
	if err != nil {
		return err
	}
	err = UpdateVideoCommentCount(&video)
	if err != nil {
		return err
	}

	// delete from redis
	err = redisService.RedisClient.HDel(redisService.Ctx, "comment", strconv.Itoa(int(comment.ID))).Err()
	return err
}

func GetCommentsByVideo(video *Video) ([]Comment, error) {
	db, _ := GetDB()
	var comments []Comment
	err := db.Where("video_id = ?", video.ID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

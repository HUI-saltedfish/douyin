package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Author        User      `json:"author"`
	AuthorID      uint      `gorm:"index"`
	PlayUrl       string    `json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int       `json:"favorite_count,omitempty" gorm:"default:0"`
	CommentCount  int       `json:"comment_count,omitempty" gorm:"default:0"`
	Is_favorite   bool      `json:"is_favorite,omitempty" gorm:"default:false"`
	Title         string    `json:"title,omitempty"`
	Favorite_User []User    `gorm:"many2many:favorite_videos;"`
	Comments      []Comment `gorm:"foreignKey:VideoID"`
}

func CreateVideo(video *Video) error {
	db, _ := GetDB()
	return db.Create(video).Error
}

func GetVideoOrderByTime() ([]Video, error) {
	db, _ := GetDB()
	var videos []Video
	err := db.Preload("Author").Order("created_at desc").Find(&videos).Error
	return videos, err
}
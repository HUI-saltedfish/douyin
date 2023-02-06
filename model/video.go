package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Author        User      `json:"author"`
	AuthorID      uint      `gorm:"index"`
	PlayUrl       string    `json:"play_url,omitempty"`
	CoverUrl      string    `json:"cover_url,omitempty"`
	FavoriteCount int       `json:"favorite_count,omitempty"`
	CommentCount  int       `json:"comment_count,omitempty"`
	Is_favorite   bool      `json:"is_favorite,omitempty"`
	Title         string    `json:"title,omitempty"`
	Favorite_User []User    `gorm:"many2many:favorite_videos;"`
	Comments      []Comment `gorm:"foreignKey:VideoID"`
}

package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content    string `gorm:"type:text;not null" json:"content,omitempty"`
	CreateDate string `gorm:"type:varchar(20);not null" json:"create_date,omitempty"`
	User       User   `gorm:"foreignKey:UserID" json:"user"`
	UserID     uint   `gorm:"not null"`
	V          Video  `gorm:"foreignKey:VideoID"`
	VideoID    uint   `gorm:"not null"`
}

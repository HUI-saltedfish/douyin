package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	CommentId         int64  `json:"id,omitempty"`
	Content    string `gorm:"type:text;not null" json:"content,omitempty"`
	CreateDate string `gorm:"type:varchar(20);not null" json:"create_date,omitempty"`
	User       User   `gorm:"foreignKey:UserID" json:"user"`
	UserID     uint   `gorm:"not null"`
	V          Video  `gorm:"foreignKey:VideoID"`
	VideoID    uint   `gorm:"not null"`
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
	return err
}

func GetCommentById(id int64) (*Comment, error) {
	db, _ := GetDB()
	var comment Comment
	err := db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	comment.CommentId = int64(comment.ID)
	return &comment, nil
}

func DeleteComment(comment *Comment) error {
	db, _ := GetDB()
	err :=  db.Delete(comment).Error
	if err != nil {
		return err
	}
	video, err := GetVideoById(uint(comment.VideoID))
	if err != nil {
		return err
	}
	err = UpdateVideoCommentCount(&video)
	return err
}

func GetCommentsByVideo(video *Video) ([]Comment, error) {
	db, _ := GetDB()
	var comments []Comment
	err := db.Where("video_id = ?", video.ID).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// update comment Id
	for i, c := range comments {
		comments[i].CommentId = int64(c.ID)
	}
	return comments, nil
}

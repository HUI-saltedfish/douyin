package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	VideoId       int64     `json:"id,omitempty" gorm:"default:0"`
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

	// update Id by gorm.ID and
	for i, v := range videos {
		videos[i].VideoId = int64(v.ID)
	}
	return videos, err
}

func GetVideoById(id uint) (Video, error) {
	db, _ := GetDB()
	var video Video
	err := db.Preload("Author").First(&video, id).Error

	// update Id by gorm.ID
	video.VideoId = int64(video.ID)
	return video, err
}

func GetVideosByUser(user *User) ([]Video, error) {
	var videos []Video
	db, _ := GetDB()
	err := db.Model(user).Association("Has_Videos").Find(&videos)

	// update Id by gorm.ID
	for i, v := range videos {
		videos[i].VideoId = int64(v.ID)
	}
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

	// update Id by gorm.ID
	for i, v := range videos {
		videos[i].VideoId = int64(v.ID)
	}
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

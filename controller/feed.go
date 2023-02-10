package controller

import (
	"log"
	"net/http"
	"simple-demo/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type FeedResponse struct {
	Response
	NextTime  int64         `json:"next_time,omitempty"`
	VideoList []model.Video `json:"video_list,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var err error
	var max_videos int = 30
	var videoList []model.Video
	if videoList, err = model.GetVideoOrderByTime(); err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	token, ok1 := c.GetQuery("token")
	if ok1 {
		username, err := getUsernameFromToken(token)
		if err == nil {
			user, err := model.GetUserByName(username)
			//没查到，一切如常
			if err == nil {
				for i := 0; i < len(videoList); i++ {
					videoList[i].Is_favorite = model.IsFavoriteVideo(user, &videoList[i])
				}
			}
		}
	}

	lastTime, ok2 := c.GetQuery("latest_time")
	if ok2 {
		log.Println("have last_time")
		last_time, err := strconv.Atoi(lastTime)
		if err != nil {
			log.Println("invalid last_time")
		}
		//保留lastTime后的视频数据
		earlistTime := -1
		for k, v := range videoList {
			if time.Unix(int64(last_time), 0).After(v.CreatedAt) {
				earlistTime = k
				break
			}
		}
		if earlistTime != -1 {
			videoList = videoList[earlistTime:]
		}
	}

	if len(videoList) > max_videos {
		videoList = videoList[:max_videos]
	}

	next_time := videoList[len(videoList)-1].CreatedAt.Unix()
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		NextTime:  next_time,
		VideoList: videoList,
	})
}

func getUsernameFromToken(token string) (username string, err error) {

	// get user by token
	// Parse the JWT string and store the result in `claims`.
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "Unauthorized access", err
	}
	return claims.Username, nil
}

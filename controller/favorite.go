package controller

import (
	"net/http"
	"simple-demo/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	// token := c.Query("token")

	// if _, exist := usersLoginInfo[token]; exist {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	user, err := model.GetUserByName(username.(string))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	video_id, ok := c.GetQuery("video_id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is empty"})
		return
	}
	vedioId, err := strconv.Atoi(video_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is not valid"})
		return
	}
	vedio, err := model.GetVideoById(uint(vedioId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is not exist"})
		return
	}

	action_type, ok := c.GetQuery("action_type")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is empty"})
		return
	}
	actionType, err := strconv.Atoi(action_type)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is not valid"})
		return
	}

	if actionType == 1 {
		err = model.AddFavoriteVideo(user, &vedio)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
	} else if actionType == 2 {
		err = model.UnFavoriteVideo(user, &vedio)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is not valid"})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	// c.JSON(http.StatusOK, VideoListResponse{
	// 	Response: Response{
	// 		StatusCode: 0,
	// 	},
	// 	VideoList: nil,
	// })
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	user, err := model.GetUserByName(username.(string))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	videos, err := model.GetUserFavoriteVideos(user)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})

}

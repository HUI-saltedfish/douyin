package controller

import (
	"net/http"
	"simple-demo/model"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	// c.JSON(http.StatusOK, FeedResponse{
	// 	Response:  Response{StatusCode: 0},
	// 	VideoList: DemoVideos,
	// 	NextTime:  time.Now().Unix(),
	// })
	var max_videos int = 30
	var videoList []model.Video
	var err error
	if videoList, err = model.GetVideoOrderByTime(); err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	if len(videoList) > max_videos {
		videoList = videoList[:max_videos]
	}

	next_time := videoList[len(videoList)-1].CreatedAt.Unix()
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  next_time,
	})
}

package controller

import (
	"context"
	"net/http"
	"simple-demo/foundationPb"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	var ret FeedResponse

	// get user by token
	last_time := c.Query("latest_time")
	last_timeInt, err := strconv.ParseInt(last_time, 10, 64)
	if err != nil {
		last_timeInt = time.Now().Unix()
	}
	if last_timeInt > time.Now().Unix() {
		last_timeInt = time.Now().Unix()
	}

	username := c.GetString("username")

	request_pd := foundationPb.DouyinFeedRequest{
		Username:   username,
		LatestTime: last_timeInt,
	}

	response_pd, err := foundationClient.Feed(context.Background(), &request_pd)
	if err != nil {
		ret.StatusCode = 1
		ret.StatusMsg = err.Error()
		ret.NextTime = time.Now().Unix()
		ret.VideoList = []Video{}
		c.JSON(http.StatusOK, ret)
		return
	}

	ret.StatusCode = response_pd.StatusCode
	ret.StatusMsg = response_pd.StatusMsg
	ret.NextTime = response_pd.NextTime
	ret.VideoList = []Video{}
	for _, video := range response_pd.VideoList {
		ret.VideoList = append(ret.VideoList, Video{
			Id: video.Id,
			Author: User{
				Id:            video.Author.Id,
				Name:          video.Author.Name,
				FollowCount:   video.Author.FollowCount,
				FollowerCount: video.Author.FollowerCount,
				IsFollow:      video.Author.IsFollow,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
		})
	}

	c.JSON(http.StatusOK, ret)
}

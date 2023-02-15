package controller

import (
	"context"
	"net/http"
	"simple-demo/interactionPb"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
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

	favActReq := interactionPb.DouyinFavoriteActionRequest{
		Username:   username.(string),
		VideoId:    int64(vedioId),
		ActionType: int32(actionType),
	}

	var favActResp *interactionPb.DouyinFavoriteActionResponse
	favActResp, err = interactionClient.FavoriteAction(context.Background(), &favActReq)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	if favActResp.StatusCode != 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: favActResp.StatusMsg})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})

	// user, err := model.GetUserByName(username.(string))
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }

	// vedio, err := model.GetVideoById(uint(vedioId))
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is not exist"})
	// 	return
	// }

	// if actionType == 1 {
	// 	err = model.AddFavoriteVideo(user, &vedio)
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	// 		return
	// 	}
	// } else if actionType == 2 {
	// 	err = model.UnFavoriteVideo(user, &vedio)
	// 	if err != nil {
	// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	// 		return
	// 	}
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is not valid"})
	// 	return
	// }

	// c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	user_id, ok := c.GetQuery("user_id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "user_id is empty"})
		return
	}
	userId, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "user_id is not valid"})
		return
	}

	favListReq := interactionPb.DouyinFavoriteListRequest{
		Username: username.(string),
		UserId:   int64(userId),
	}

	var favListResp *interactionPb.DouyinFavoriteListResponse
	favListResp, err = interactionClient.FavoriteList(context.Background(), &favListReq)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
		return
	}

	if favListResp.StatusCode != 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: favListResp.StatusMsg})
		return
	}

	var videos []Video
	for _, v := range favListResp.VideoList {
		videos = append(videos, Video{
			Id: v.Id,
			Author: User{
				Id:            v.Author.Id,
				Name:          v.Author.Name,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
				IsFollow:      v.Author.IsFollow,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
		})
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})

}

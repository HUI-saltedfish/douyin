package controller

import (
	"context"
	"net/http"
	"simple-demo/interactionPb"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	// token := c.Query("token")
	// actionType := c.Query("action_type")

	// if user, exist := usersLoginInfo[token]; exist {
	// 	if actionType == "1" {
	// 		text := c.Query("comment_text")
	// 		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
	// 			Comment: Comment{
	// 				Id:         1,
	// 				User:       user,
	// 				Content:    text,
	// 				CreateDate: "05-01",
	// 			}})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
	// } else {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// }

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
	videoId, err := strconv.Atoi(video_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is invalid"})
		return
	}

	action_type, ok := c.GetQuery("action_type")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is empty"})
		return
	}
	actionType, err := strconv.Atoi(action_type)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is invalid"})
		return
	}
	comment_text := c.Query("comment_text")
	comment_id := c.Query("comment_id")
	comment_id_int, err := strconv.ParseInt(comment_id, 10, 64)
	if err != nil && actionType == 2 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_id is invalid"})
		return
	}

	commActReq := &interactionPb.DouyinCommentActionRequest{
		Username:    username.(string),
		VideoId:     int64(videoId),
		ActionType:  int32(actionType),
		CommentText: comment_text,
		CommentId:   comment_id_int,
	}
	commActResp, err := interactionClient.CommentAction(context.Background(), commActReq)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment action failed"})
		return
	}
	if commActResp.StatusCode != 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: commActResp.StatusMsg})
		return
	}

	if actionType == 1 {
		ret := Comment{
			Id: commActResp.Comment.Id,
			User: User{
				Id:            commActResp.Comment.User.Id,
				Name:          commActResp.Comment.User.Name,
				FollowCount:   commActResp.Comment.User.FollowCount,
				FollowerCount: commActResp.Comment.User.FollowerCount,
				IsFollow:      commActResp.Comment.User.IsFollow,
			},
			Content:    commActResp.Comment.Content,
			CreateDate: commActResp.Comment.CreateDate,
		}
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0}, Comment: ret})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// if actionType == 1 {
// 	comment_text, ok := c.GetQuery("comment_text")
// 	if !ok {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_text is empty"})
// 		return
// 	}
// 	comment := model.Comment{
// 		Content:    comment_text,
// 		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
// 		UserID:     user.ID,
// 		VideoID:    video.ID,
// 	}
// 	err := model.CreateComment(&comment)
// 	if err != nil {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "create comment failed"})
// 		return
// 	}
// 	c.JSON(
// 		http.StatusOK,
// 		CommentActionResponse{
// 			Response: Response{StatusCode: 0},
// 			Comment:  comment,
// 		},
// 	)
// } else if actionType == 2 {
// 	comment_id, ok := c.GetQuery("comment_id")
// 	if !ok {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_id is empty"})
// 		return
// 	}
// 	commentId, err := strconv.Atoi(comment_id)
// 	if err != nil {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_id is invalid"})
// 		return
// 	}
// 	comment, err := model.GetCommentById(int64(commentId))
// 	if err != nil {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment doesn't exist"})
// 		return
// 	}
// 	err = model.DeleteComment(comment)
// 	if err != nil {
// 		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "delete comment failed"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, Response{StatusCode: 0})
// } else {
// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is invalid"})
// 	return
// }

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	video_id, ok := c.GetQuery("video_id")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is empty"})
		return
	}
	videoId, err := strconv.Atoi(video_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is invalid"})
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	comListReq := interactionPb.DouyinCommentListRequest{
		Username: username.(string),
		VideoId:  int64(videoId),
	}
	comListRes, err := interactionClient.CommentList(context.Background(), &comListReq)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get comments failed"})
		return
	}
	if comListRes.StatusCode != 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: comListRes.StatusMsg})
		return
	}

	comments := make([]Comment, 0)
	for _, com := range comListRes.CommentList {
		comments = append(comments, Comment{
			Id: com.Id,
			User: User{
				Id:            com.User.Id,
				Name:          com.User.Name,
				FollowCount:   com.User.FollowCount,
				FollowerCount: com.User.FollowerCount,
				IsFollow:      com.User.IsFollow,
			},
			Content:    com.Content,
			CreateDate: com.CreateDate,
		})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}

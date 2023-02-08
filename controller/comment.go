package controller

import (
	"net/http"
	"simple-demo/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListResponse struct {
	Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment model.Comment `json:"comment,omitempty"`
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
	videoId, err := strconv.Atoi(video_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video_id is invalid"})
		return
	}
	video, err := model.GetVideoById(uint(videoId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video doesn't exist"})
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

	if actionType == 1 {
		comment_text, ok := c.GetQuery("comment_text")
		if !ok {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_text is empty"})
			return
		}
		comment := model.Comment{
			Content:    comment_text,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
			UserID:     user.ID,
			VideoID:    video.ID,
		}
		err := model.CreateComment(&comment)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "create comment failed"})
			return
		}
		c.JSON(
			http.StatusOK,
			CommentActionResponse{
				Response: Response{StatusCode: 0},
				Comment:  comment,
			},
		)
	} else if actionType == 2 {
		comment_id, ok := c.GetQuery("comment_id")
		if !ok {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_id is empty"})
			return
		}
		commentId, err := strconv.Atoi(comment_id)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment_id is invalid"})
			return
		}
		comment, err := model.GetCommentById(int64(commentId))
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "comment doesn't exist"})
			return
		}
		// delete the association between comment and user
		// err = model.DeleteCommentAssociationByUser(user, comment)
		// if err != nil {
		// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "delete comment and user association failed"})
		// 	return
		// }
		// delete the association between comment and video
		// err = model.DeleteCommentAssociationByVideo(&video, comment)
		// if err != nil {
		// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "delete comment and video association failed"})
		// 	return
		// }
		err = model.DeleteComment(comment)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "delete comment failed"})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "action_type is invalid"})
		return
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	// c.JSON(http.StatusOK, CommentListResponse{
	// 	Response:    Response{StatusCode: 0},
	// 	CommentList: DemoComments,
	// })
	// username, ok := c.Get("username")
	// if !ok {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }

	// user, err := model.GetUserByName(username.(string))
	// if err != nil {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }

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
	video, err := model.GetVideoById(uint(videoId))
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "video doesn't exist"})
		return
	}
	comments, err := model.GetCommentsByVideo(&video)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get comments failed"})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}

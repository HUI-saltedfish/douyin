package controller

import (
	"context"
	"mime/multipart"
	"net/http"
	"simple-demo/foundationPb"
	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	// getuser
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	titel, ok := c.GetPostForm("title")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "title is empty",
		})
		return
	}
	video_url, err := UploadAliyunOss(data)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	pubActReq := foundationPb.DouyinPublishActionRequest{
		Username: username.(string),
		PlayUrl:  video_url,
		Title:    titel,
	}

	pubActRes, err := foundationClient.PublishAction(context.Background(), &pubActReq)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: pubActRes.StatusCode,
		StatusMsg:  data.Filename + " uploaded successfully",
	})

}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}
	user_id, ok := c.GetQuery("user_id")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "user_id is empty",
		})
		return
	}
	userID, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	publishListRequest := foundationPb.DouyinPublishListRequest{
		Username: username.(string),
		UserId:   int64(userID),
	}

	publishListResponse, err := foundationClient.PublishList(context.Background(), &publishListRequest)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	var pubListRes VideoListResponse
	pubListRes.StatusCode = publishListResponse.StatusCode
	for _, v := range publishListResponse.VideoList {
		pubListRes.VideoList = append(pubListRes.VideoList, Video{
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

	c.JSON(http.StatusOK, pubListRes)
}

// 后期可以改微服务
func UploadAliyunOss(file *multipart.FileHeader) (url string, err error) {
	Endpoint := "oss-cn-beijing.aliyuncs.com"
	AccessKeyId := "LTAI5tGqn5BjBBk9ecXLuGzU"
	AccessKeySecret := "Z4dkbCnUjCnqdQbCBFPQXlmQYtOfoE"
	BucketName := "myblog-hui"

	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		panic(err)
	}

	// 指定自己的bucket
	bucket, err := client.Bucket(BucketName)
	if err != nil {
		panic(err)
	}

	src, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	// 上传文件 并返回文件的URL
	path := "video/" + file.Filename
	err = bucket.PutObject(path, src)
	if err != nil {
		panic(err)
	}
	return "https://" + BucketName + "." + Endpoint + "/" + path, nil
}

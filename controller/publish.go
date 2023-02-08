package controller

import (
	"mime/multipart"
	"net/http"
	"simple-demo/model"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	// token := c.PostForm("token")

	// if _, exist := usersLoginInfo[token]; !exist {
	// 	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	// 	return
	// }

	// getuser
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}
	user, err := model.GetUserByName(username.(string))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
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

	// filename := filepath.Base(data.Filename)
	// user := usersLoginInfo[token]
	// finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	// saveFile := filepath.Join("./public/", finalName)
	// if err := c.SaveUploadedFile(data, saveFile); err != nil {
	// 	c.JSON(http.StatusOK, Response{
	// 		StatusCode: 1,
	// 		StatusMsg:  err.Error(),
	// 	})
	// 	return
	// }
	titel, ok := c.GetPostForm("title")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "title is empty",
		})
		return
	}
	vedio_url, err := UploadAliyunOss(data)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	vedio := model.Video{
		AuthorID:      user.ID,
		PlayUrl:       vedio_url,
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		Is_favorite:   false,
		Title:         titel,
	}
	err = model.CreateVideo(&vedio)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  data.Filename + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// c.JSON(http.StatusOK, VideoListResponse{
	// 	Response: Response{
	// 		StatusCode: 0,
	// 	},
	// 	VideoList: DemoVideos,
	// })
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User doesn't exist",
		})
		return
	}
	user, err := model.GetUserByName(username.(string))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videos, err := model.GetVideosByUser(user)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// update video's is_favorite field
	for i := 0; i < len(videos); i++ {
		videos[i].Is_favorite = model.IsFavoriteVideo(user, &videos[i])
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
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

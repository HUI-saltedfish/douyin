package main

import (
	"context"
	Pb "foundation/foundationPb"
	"foundation/model"
	"foundation/redisService"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type foundationServer struct {
	Pb.UnimplementedDouyinFoundationServiceServer
	savedResults []*Pb.DouyinFeedResponse
}

func (s *foundationServer) Feed(ctx context.Context, in *Pb.DouyinFeedRequest) (*Pb.DouyinFeedResponse, error) {
	lastTime := in.LatestTime
	username := in.Username

	if lastTime == 0 {
		lastTime = time.Now().Unix()
	}
	// convert to time.Time
	t := time.Unix(lastTime, 0)

	// get videos
	videos, err := model.GetVideoOrderByTime(t)
	if err != nil {
		log.Println("GetVideoOrderByTime error: ", err)
		return &Pb.DouyinFeedResponse{
			StatusCode: 1,
			StatusMsg:  "GetVideoOrderByTime error",
		}, nil
	}

	if len(videos) == 0 {
		log.Println("this is the last video")
		return &Pb.DouyinFeedResponse{
			StatusCode: 1,
			StatusMsg:  "this is the last video",
		}, nil
	}

	// judge if user has login
	var user *model.User
	if username != "" {
		// get user
		user, err = model.GetUserByName(username)
		if err != nil {
			log.Println("GetUserByName error: ", err)
			return &Pb.DouyinFeedResponse{
				StatusCode: 1,
				StatusMsg:  "GetUserByName error",
			}, nil
		}
	}

	// judge the user favorite status
	for i, v := range videos {
		var tempUser *model.User
		tempUser, _ = model.GetUserById(int(v.AuthorID))
		videos[i].Author = *tempUser
		if username != "" {
			// judge if user has liked the video
			videos[i].Is_favorite = model.IsFavoriteVideo(user, &v)
		} else {
			videos[i].Is_favorite = false
		}
	}

	// convert to pb
	pbVideos := make([]*Pb.Video, len(videos))
	for i, v := range videos {
		pbVideos[i] = &Pb.Video{
			Id: int64(v.ID),
			Author: &Pb.User{
				Id:            int64(v.Author.ID),
				Name:          v.Author.Name,
				IsFollow:      v.Author.Is_follow,
				FollowCount:   int64(v.Author.Follow_count),
				FollowerCount: int64(v.Author.Follower_count),
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    v.Is_favorite,
			Title:         v.Title,
		}
	}

	// save result
	s.savedResults = append(s.savedResults, &Pb.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  pbVideos,
		NextTime:   videos[0].CreatedAt.Unix(),
	})

	return &Pb.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  pbVideos,
		NextTime:   videos[0].CreatedAt.Unix(),
	}, nil

}

func (s *foundationServer) UserRegister(ctx context.Context, in *Pb.DouyinUserRegisterRequest) (*Pb.DouyinUserRegisterResponse, error) {
	username := in.Username
	password := in.Password

	// check if username exists
	var user *model.User
	var err error
	_, err = model.GetUserByName(username)
	if err == nil {
		log.Println("username exists")
		return &Pb.DouyinUserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  "username exists",
		}, nil
	}

	// create user
	user = &model.User{
		Name:     username,
		Password: password,
	}
	err = model.CreateUser(user)
	if err != nil {
		log.Println("CreateUser error: ", err)
		return &Pb.DouyinUserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  "CreateUser error",
		}, nil
	}

	return &Pb.DouyinUserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(user.ID),
		Username:   user.Name,
	}, nil
}

func (s *foundationServer) UserLogin(ctx context.Context, in *Pb.DouyinUserLoginRequest) (*Pb.DouyinUserLoginResponse, error) {
	username := in.Username
	password := in.Password

	// get user
	user, err := model.GetUserByName(username)
	if err != nil {
		log.Println("GetUserByName error: ", err)
		return &Pb.DouyinUserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserByName error",
		}, nil
	}

	// check password
	if user.Password != password {
		log.Println("password error")
		return &Pb.DouyinUserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "password error",
		}, nil
	}

	return &Pb.DouyinUserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     int64(user.ID),
		Username:   user.Name,
	}, nil
}

func (s *foundationServer) User(ctx context.Context, in *Pb.DouyinUserRequest) (*Pb.DouyinUserResponse, error) {
	userId := in.UserId

	// get user
	var user *model.User
	var err error
	user, err = model.GetUserById(int(userId))
	if err != nil {
		log.Println("GetUserById error: ", err)
		return &Pb.DouyinUserResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserById error",
		}, nil
	}

	return &Pb.DouyinUserResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		User: &Pb.User{
			Id:            int64(user.ID),
			Name:          user.Name,
			IsFollow:      user.Is_follow,
			FollowCount:   int64(user.Follow_count),
			FollowerCount: int64(user.Follower_count),
		},
	}, nil

}

func (s *foundationServer) PublishList(ctx context.Context, in *Pb.DouyinPublishListRequest) (*Pb.DouyinPublishListResponse, error) {
	userId := in.UserId

	// get user
	var user *model.User
	var err error
	user, err = model.GetUserById(int(userId))
	if err != nil {
		log.Println("GetUserById error: ", err)
		return &Pb.DouyinPublishListResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserById error",
		}, nil
	}

	// get videos
	var videos []model.Video
	videos, err = model.GetVideosByUser(user)
	if err != nil {
		log.Println("GetVideosByAuthor error: ", err)
		return &Pb.DouyinPublishListResponse{
			StatusCode: 1,
			StatusMsg:  "GetVideosByAuthor error",
		}, nil
	}

	// judge the user favorite status
	for i, v := range videos {
		videos[i].Is_favorite = model.IsFavoriteVideo(user, &v)
	}

	// convert to pb
	pbVideos := make([]*Pb.Video, len(videos))
	for i, v := range videos {
		pbVideos[i] = &Pb.Video{
			Id: int64(v.ID),
			Author: &Pb.User{
				Id:            int64(v.Author.ID),
				Name:          v.Author.Name,
				FollowCount:   int64(v.Author.Follow_count),
				FollowerCount: int64(v.Author.Follower_count),
				IsFollow:      v.Author.Is_follow,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    v.Is_favorite,
			Title:         v.Title,
		}
	}

	return &Pb.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  pbVideos,
	}, nil
}

func (s *foundationServer) PublishAction(ctx context.Context, in *Pb.DouyinPublishActionRequest) (*Pb.DouyinPublishActionResponse, error) {
	username := in.Username
	video_url := in.PlayUrl
	title := in.Title

	// get user
	var user *model.User
	var err error
	user, err = model.GetUserByName(username)
	if err != nil {
		log.Println("GetUserByName error: ", err)
		return &Pb.DouyinPublishActionResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserByName error",
		}, nil
	}

	// create video
	var video = model.Video{
		AuthorID:      user.ID,
		PlayUrl:       video_url,
		CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		Is_favorite:   false,
		Title:         title,
	}

	err = model.CreateVideo(&video)
	if err != nil {
		log.Println("CreateVideo error: ", err)
		return &Pb.DouyinPublishActionResponse{
			StatusCode: 1,
			StatusMsg:  "CreateVideo error",
		}, nil
	}

	return &Pb.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

func main() {
	// init
	model.GetDB()
	defer redisService.RedisClient.Close()

	// start server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	Pb.RegisterDouyinFoundationServiceServer(s, &foundationServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

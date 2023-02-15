package main

import (
	"context"
	"interaction/interactionPb"
	"interaction/model"
	"interaction/redisService"
	"log"
	"net"

	"google.golang.org/grpc"
)

type interactionServer struct {
	interactionPb.UnimplementedDouyinInteractionServiceServer
}

func (s *interactionServer) FavoriteAction(ctx context.Context, in *interactionPb.DouyinFavoriteActionRequest) (*interactionPb.DouyinFavoriteActionResponse, error) {
	username := in.GetUsername()
	videoId := in.GetVideoId()
	action_type := in.GetActionType()

	user, err := model.GetUserByName(username)
	if err != nil {
		log.Println("GetUserByName error: ", err)
		return &interactionPb.DouyinFavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserByName error",
		}, nil
	}

	video, err := model.GetVideoById(uint(videoId))
	if err != nil {
		log.Println("GetVideoById error: ", err)
		return &interactionPb.DouyinFavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  "GetVideoById error",
		}, nil
	}

	switch action_type {
	case 1:
		err = model.AddFavoriteVideo(user, &video)
	case 2:
		err = model.UnFavoriteVideo(user, &video)
	default:
		log.Println("action_type error: ", action_type)
		return &interactionPb.DouyinFavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  "action_type error",
		}, nil
	}
	if err != nil {
		log.Println("FavoriteAction error: ", err)
		return &interactionPb.DouyinFavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		}, nil
	}

	return &interactionPb.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}

func (s *interactionServer) FavoriteList(ctx context.Context, in *interactionPb.DouyinFavoriteListRequest) (*interactionPb.DouyinFavoriteListResponse, error) {
	userId := in.GetUserId()

	user, err := model.GetUserById(int(userId))
	if err != nil {
		log.Println("GetUserById error: ", err)
		return &interactionPb.DouyinFavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserById error",
		}, nil
	}

	favoriteVideos, err := model.GetUserFavoriteVideos(user)
	if err != nil {
		log.Println("GetUserFavoriteVideos error: ", err)
		return &interactionPb.DouyinFavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserFavoriteVideos error",
		}, nil
	}

	// update the user favorite status
	for i := 0; i < len(favoriteVideos); i++ {
		favoriteVideos[i].Is_favorite = model.IsFavoriteVideo(user, &favoriteVideos[i])
	}

	var videos []*interactionPb.Video
	for _, video := range favoriteVideos {
		videos = append(videos, &interactionPb.Video{
			Id: int64(video.ID),
			Author: &interactionPb.User{
				Id:            int64(video.Author.ID),
				Name:          video.Author.Name,
				FollowCount:   int64(video.Author.Follow_count),
				FollowerCount: int64(video.Author.Follower_count),
				IsFollow:      video.Is_favorite,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
			IsFavorite:    video.Is_favorite,
			Title:         video.Title,
		})
	}

	return &interactionPb.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videos,
	}, nil

}

func (s *interactionServer) CommentList(ctx context.Context, in *interactionPb.DouyinCommentListRequest) (*interactionPb.DouyinCommentListResponse, error) {
	videoId := in.GetVideoId()

	video, err := model.GetVideoById(uint(videoId))
	if err != nil {
		log.Println("GetVideoById error: ", err)
		return &interactionPb.DouyinCommentListResponse{
			StatusCode: 1,
			StatusMsg:  "GetVideoById error",
		}, nil
	}

	comments, err := model.GetCommentsByVideo(&video)
	if err != nil {
		log.Println("GetVideoComments error: ", err)
		return &interactionPb.DouyinCommentListResponse{
			StatusCode: 1,
			StatusMsg:  "GetVideoComments error",
		}, nil
	}

	// get user and video info
	for _, comment := range comments {
		tempUser, err := model.GetUserById(int(comment.UserID))
		if err != nil {
			log.Println("GetUserById error: ", err)
			return &interactionPb.DouyinCommentListResponse{
				StatusCode: 1,
				StatusMsg:  "GetUserById error",
			}, nil
		}
		comment.User = *tempUser
	}

	var commentList []*interactionPb.Comment
	for _, comment := range comments {
		commentList = append(commentList, &interactionPb.Comment{
			Id: int64(comment.ID),
			User: &interactionPb.User{
				Id:            int64(comment.User.ID),
				Name:          comment.User.Name,
				FollowCount:   int64(comment.User.Follow_count),
				FollowerCount: int64(comment.User.Follower_count),
				IsFollow:      comment.User.Is_follow,
			},
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &interactionPb.DouyinCommentListResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		CommentList: commentList,
	}, nil

}

func (s *interactionServer) CommentAction(ctx context.Context, in *interactionPb.DouyinCommentActionRequest) (*interactionPb.DouyinCommentActionResponse, error) {
	username := in.GetUsername()
	videoId := in.GetVideoId()
	action_type := in.GetActionType()
	content := in.GetCommentText()
	comment_id := in.GetCommentId()

	user, err := model.GetUserByName(username)
	if err != nil {
		log.Println("GetUserById error: ", err)
		return &interactionPb.DouyinCommentActionResponse{
			StatusCode: 1,
			StatusMsg:  "GetUserById error",
		}, nil
	}

	video, err := model.GetVideoById(uint(videoId))
	if err != nil {
		log.Println("GetVideoById error: ", err)
		return &interactionPb.DouyinCommentActionResponse{
			StatusCode: 1,
			StatusMsg:  "GetVideoById error",
		}, nil
	}

	var comm model.Comment
	switch action_type {
	case 1:
		comm = model.Comment{
			Content: content,
			User:    *user,
			UserID:  user.ID,
			V:       video,
			VideoID: video.ID,
		}
		err = model.CreateComment(&comm)
		if err != nil {
			log.Println("CreateComment error: ", err)
			return &interactionPb.DouyinCommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "CreateComment error",
			}, nil
		}
	case 2:
		comm, err := model.GetCommentById(comment_id)
		if err != nil {
			log.Println("GetCommentById error: ", err)
			return &interactionPb.DouyinCommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "GetCommentById error",
			}, nil
		}
		err = model.DeleteComment(comm)
		if err != nil {
			log.Println("DeleteComment error: ", err)
			return &interactionPb.DouyinCommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "DeleteComment error",
			}, nil
		}
	default:
		log.Println("action_type error: ", action_type)
		return &interactionPb.DouyinCommentActionResponse{
			StatusCode: 1,
			StatusMsg:  "action_type error",
		}, nil
	}

	if action_type == 1 {
		return &interactionPb.DouyinCommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "success",
			Comment: &interactionPb.Comment{
				Id: int64(comm.ID),
				User: &interactionPb.User{
					Id:            int64(comm.User.ID),
					Name:          comm.User.Name,
					FollowCount:   int64(comm.User.Follow_count),
					FollowerCount: int64(comm.User.Follower_count),
					IsFollow:      comm.User.Is_follow,
				},
				Content:    comm.Content,
				CreateDate: comm.CreatedAt.Format("2006-01-02 15:04:05"),
			},
		}, nil
	} else {
		return &interactionPb.DouyinCommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		}, nil
	}
}

func main() {
	// init DB
	model.GetDB()
	defer redisService.RedisClient.Close()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	interactionPb.RegisterDouyinInteractionServiceServer(s, &interactionServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

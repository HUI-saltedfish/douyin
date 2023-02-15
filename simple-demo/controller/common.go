package controller

import (
	"simple-demo/foundationPb"
	"simple-demo/interactionPb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

var (
	foundationPbAddress = "localhost:50051"
	interactionPbAddress = "localhost:50052"
	FoundationPbConn    *grpc.ClientConn
	InteractionPbConn   *grpc.ClientConn
	foundationClient    foundationPb.DouyinFoundationServiceClient
	interactionClient   interactionPb.DouyinInteractionServiceClient
)

func init() {
	var err error

	FoundationPbConn, err := grpc.Dial(foundationPbAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	InteractionPbConn, err := grpc.Dial(interactionPbAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	foundationClient = foundationPb.NewDouyinFoundationServiceClient(FoundationPbConn)
	interactionClient = interactionPb.NewDouyinInteractionServiceClient(InteractionPbConn)

}

package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/app/api/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/comment"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc user.User

	VideoRpc video.Video

	FollowRpc follow.Follow

	FavoriteRpc favorite.Favorite

	CommentRpc comment.Comment
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		UserRpc:     user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:    video.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		FollowRpc:   follow.NewFollow(zrpc.MustNewClient(c.FollowRpc)),
		FavoriteRpc: favorite.NewFavorite(zrpc.MustNewClient(c.FavoriteRpc)),
		CommentRpc:  comment.NewComment(zrpc.MustNewClient(c.CommentRpc)),
	}
}

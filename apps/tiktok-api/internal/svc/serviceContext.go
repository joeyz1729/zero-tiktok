package svc

import (
	"github.com/joeyz1729/zero-tiktok/apps/comment/rpc/comment"
	"github.com/joeyz1729/zero-tiktok/apps/favorite/rpc/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/follow/rpc/follow"
	"github.com/joeyz1729/zero-tiktok/apps/message/message"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRpc     userservice.UserService
	VideoRpc    videoservice.VideoService
	FollowRpc   follow.Follow
	FavoriteRpc favorite.Favorite
	CommentRpc  comment.Comment
	MessageRpc  message.Message
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		UserRpc:     userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:    videoservice.NewVideoService(zrpc.MustNewClient(c.VideoRpc)),
		FollowRpc:   follow.NewFollow(zrpc.MustNewClient(c.FollowRpc)),
		FavoriteRpc: favorite.NewFavorite(zrpc.MustNewClient(c.FavoriteRpc)),
		CommentRpc:  comment.NewComment(zrpc.MustNewClient(c.CommentRpc)),
		MessageRpc:  message.NewMessage(zrpc.MustNewClient(c.MessageRpc)),
	}
}

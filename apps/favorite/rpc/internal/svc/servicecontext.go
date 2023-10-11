package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/dao"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	FavorRepo *dao.RepoImpl

	UserRpc user.User

	VideoRpc video.Video
}

func NewServiceContext(c config.Config) *ServiceContext {
	r, err := dao.NewRepo(c)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		FavorRepo: r,
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:  video.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
	}
}

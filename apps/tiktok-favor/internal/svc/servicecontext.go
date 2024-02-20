package svc

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/internal/repository"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/videoservice"

	"github.com/zeromicro/go-zero/zrpc"
	_ "gorm.io/driver/mysql"
)

type ServiceContext struct {
	Config config.Config

	ThumbupDao *repository.ThumbupDao

	UserRpc userservice.UserService

	VideoRpc videoservice.VideoService
}

func NewServiceContext(c config.Config) *ServiceContext {
	r, err := repository.NewRepo(c.Mysql.DataSource, c.CacheRedis.Addr)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:     c,
		ThumbupDao: r,
		UserRpc:    userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:   videoservice.NewVideoService(zrpc.MustNewClient(c.VideoRpc)),
	}
}

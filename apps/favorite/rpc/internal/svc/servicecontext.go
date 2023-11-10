package svc

import (
	"database/sql"

	"github.com/YiZou89/zero-tiktok/apps/favorite/data"
	"github.com/YiZou89/zero-tiktok/apps/favorite/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/user"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/video"

	"github.com/zeromicro/go-zero/zrpc"
	_ "gorm.io/driver/mysql"
)

type ServiceContext struct {
	Config config.Config

	FavorRepo *data.RepoImpl

	UserRpc user.User

	VideoRpc video.Video

	BarrierDB *sql.DB

	DtmServer string
}

func NewServiceContext(c config.Config) *ServiceContext {
	r, err := data.NewRepo(c.Mysql.DataSource, c.CacheRedis.Addr)
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		FavorRepo: r,
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpc)),
		VideoRpc:  video.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		BarrierDB: db,
	}
}

package svc

import (
	"github.com/YiZou89/zero-tiktok/apps/follow/rpc/follow"
	"github.com/YiZou89/zero-tiktok/apps/user/data"
	"github.com/YiZou89/zero-tiktok/apps/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	//UserModel model2.UserModel

	FollowRpc follow.Follow

	UserRepo *data.Repo
}

func NewServiceContext(c config.Config) *ServiceContext {
	//sqlConn := sqlx.NewMysql(c.Repo.DataSource)
	return &ServiceContext{
		Config: c,
		//UserModel: model2.NewUserModel(sqlConn),
		FollowRpc: follow.NewFollow(zrpc.MustNewClient(c.FollowRpc)),
		UserRepo:  data.NewRepo(c.Repo.DataSource, c.Repo.RedisAddr),
	}
}

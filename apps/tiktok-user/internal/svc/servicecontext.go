package svc

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-relation/follow"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/repository"
	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	UserRepo    *repository.Repo
	IDGenerator *snowflake.IDGenerator

	FollowRpc follow.Follow
}

func NewServiceContext(c config.Config) *ServiceContext {
	repo, err := repository.NewRepo(c.Repo)
	if err != nil {
		panic(err)
	}
	idGenerator, err := snowflake.NewIDGenerator(c.Snowflake.StartTime, c.Snowflake.MachineId)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:      c,
		UserRepo:    repo,
		IDGenerator: idGenerator,
		FollowRpc:   follow.NewFollow(zrpc.MustNewClient(c.FollowRpc)),
	}
}

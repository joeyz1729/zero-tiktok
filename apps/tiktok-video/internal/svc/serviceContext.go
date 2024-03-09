package svc

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-favor/favorite"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/userservice"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/job"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/repository"
	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	IfUploadMinIO bool
	Repo          *repository.Repo
	Worker        *job.Worker
	IDGenerator   *snowflake.IDGenerator
	UserRpc       userservice.UserService
	FavorRpc      favorite.Favorite
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

	readerConf := kafka.ReaderConfig{
		Brokers:     c.Kafka.Brokers,
		Partition:   0,
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	}
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:        c,
		IfUploadMinIO: c.MinIO.Upload,
		Repo:          repo,
		Worker: &job.Worker{
			Repo:         repo,
			ReaderConfig: readerConf,
		},
		IDGenerator: idGenerator,
		UserRpc:     userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		FavorRpc:    favorite.NewFavorite(zrpc.MustNewClient(c.FavorRpc)),
	}
}

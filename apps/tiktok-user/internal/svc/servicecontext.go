package svc

import (
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-user/internal/job"
	kafka "github.com/segmentio/kafka-go"

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

	Worker *job.Worker

	FollowRpc follow.Follow
}

func NewServiceContext(c config.Config) *ServiceContext {
	repo, err := repository.NewRepo(c.Repo.DataSource, c.Repo.RedisAddr, c.Repo.EsAddresses)
	if err != nil {
		panic(err)
	}
	idGenerator, err := snowflake.NewIDGenerator(c.Snowflake.StartTime, c.Snowflake.MachineId)
	if err != nil {
		panic(err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Kafka.Brokers,
		Topic:       c.Kafka.Topic,
		Partition:   0,
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})

	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	return &ServiceContext{
		Config:      c,
		UserRepo:    repo,
		IDGenerator: idGenerator,
		Worker: &job.Worker{
			Repo:        repo,
			KafkaReader: reader,
		},
		FollowRpc: follow.NewFollow(zrpc.MustNewClient(c.FollowRpc)),
	}
}

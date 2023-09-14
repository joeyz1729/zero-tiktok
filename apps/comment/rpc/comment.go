package main

import (
	"flag"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/kmq"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/server"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/comment/rpc/model"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/comment.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		model.RegisterCommentServer(grpcServer, server.NewCommentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	_ = consul.RegisterService(c.ListenOn, c.Consul)

	if err := snowflake.Init(c.Snowflake.StartTime, c.Snowflake.MachineId); err != nil {
		panic(err)
	}

	go func() {
		q := kmq.NewMq(c.KafkaMq, ctx.CommentDB)
		queue := kq.MustNewQueue(c.KafkaMq, kq.WithHandle(q.Consume))
		defer queue.Stop()
		fmt.Println("Starting kafka mq server at 19092")
		queue.Start()
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()

}

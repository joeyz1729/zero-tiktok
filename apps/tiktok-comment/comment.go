package main

import (
	"flag"
	"fmt"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/server"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-comment/pb"
	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"

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
		pb.RegisterCommentServer(grpcServer, server.NewCommentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	if err := snowflake.Init(c.Snowflake.StartTime, c.Snowflake.MachineId); err != nil {
		panic(err)
	}

	//go initMq(c, ctx)

	fmt.Printf("Starting tiktok-user server at %s...\n", c.ListenOn)
	s.Start()
}

//func initMq(c config.Config, ctx *svc.ServiceContext) {
//	q := kmq.NewMq(c.KafkaMq, ctx.CommentDB, ctx.CommentCache)
//	queue := kq.MustNewQueue(c.KafkaMq, kq.WithHandle(q.Consume))
//	defer queue.Stop()
//	fmt.Printf("Starting kafka mq server at %v", c.KafkaMq.Brokers)
//	queue.Start()
//}

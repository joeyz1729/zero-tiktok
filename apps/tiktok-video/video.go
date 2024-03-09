package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/server"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/internal/svc"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-video/pb"
	"github.com/joeyz1729/zero-tiktok/pkg/snowflake"
	"github.com/joeyz1729/zero-tiktok/pkg/tool"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/video.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterVideoServiceServer(grpcServer, server.NewVideoServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	err := snowflake.Init(c.Snowflake.StartTime, c.Snowflake.MachineId)
	if err != nil {
		panic("snowflake initialization failed")
	}

	tool.NewSalt(c.Salt)

	go ctx.Worker.CreateVideoStart(context.TODO())
	fmt.Printf("Starting tiktok-user server at %s...\n", c.ListenOn)
	s.Start()
}

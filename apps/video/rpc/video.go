package main

import (
	"flag"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/server"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/internal/svc"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/model"
	"github.com/YiZou89/zero-tiktok/apps/video/rpc/mw/minio"
	"github.com/YiZou89/zero-tiktok/pkg/snowflake"
	"github.com/YiZou89/zero-tiktok/pkg/tool"

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
		model.RegisterVideoServer(grpcServer, server.NewVideoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	//_ = consul.RegisterService(c.ListenOn, c.Consul)

	err := snowflake.Init(c.Snowflake.StartTime, c.Snowflake.MachineId)
	if err != nil {
		panic("snowflake initialization failed")
	}

	minio.Init()

	tool.NewSalt(c.Salt)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

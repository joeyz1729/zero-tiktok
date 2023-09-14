package main

import (
	"flag"
	"fmt"
	"github.com/YiZou89/zero-tiktok/apps/comment/rmq/internal/config"
	"github.com/YiZou89/zero-tiktok/apps/comment/rmq/internal/handler"
	"github.com/YiZou89/zero-tiktok/apps/comment/rmq/internal/service"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/comment-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := service.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

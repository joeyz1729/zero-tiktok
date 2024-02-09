package main

import (
	"flag"
	"fmt"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/config"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/handler"
	"github.com/joeyz1729/zero-tiktok/apps/tiktok-api/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"net/http"
)

var configFile = flag.String("f", "etc/tiktok-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next(w, r)
		}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

package main

import (
	"github.com/gin-gonic/gin"
	"XConf/agent-api/config"
	"XConf/agent-api/handler"
	pconfig "XConf/proto/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func main() {
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		// etcd的地址：http://127.0.0.1:2379
		op.Addrs = []string{"http://127.0.0.1:2379"}
	})
	service := web.NewService(
		web.Name("go.micro.api.agent"),
		web.Registry(reg),
	)

	if err := service.Init(); err != nil {
		panic(err)
	}

	client := pconfig.NewConfigService("go.micro.srv.config", service.Options().Service.Client())

	config.Init(client, 1024*1024)
	router := Router()
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

func Router() *gin.Engine {
	router := gin.Default()
	r := router.Group("/agent/api/v1")
	r.GET("/config", handler.ReadConfig)
	r.GET("/watch", handler.WatchUpdate)

	return router
}

package main

import (
	"github.com/gin-gonic/gin"
	"XConf/admin-api/config"
	"XConf/admin-api/handler"
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
		web.Name("go.micro.api.admin"),
		web.Registry(reg),
	)

	if err := service.Init(); err != nil {
		panic(err)
	}

	client := pconfig.NewConfigService("go.micro.srv.config", service.Options().Service.Client())

	config.Init(client)
	router := Router()
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		panic(err)
	}
}

func Router() *gin.Engine {
	router := gin.Default()
	r := router.Group("/admin/api/v1")
	r.GET("/apps", handler.ListApps)
	r.POST("/app", handler.CreateApp)
	r.DELETE("/app", handler.DeleteApp)

	r.GET("/clusters", handler.ListClusters)
	r.POST("/cluster", handler.CreateCluster)
	r.DELETE("/cluster", handler.DeleteCluster)

	r.GET("/namespaces", handler.ListNamespaces)
	r.POST("/namespace", handler.CreateNamespace)
	r.DELETE("/namespace", handler.DeleteNamespace)

	r.POST("/config", handler.UpdateConfig)
	r.POST("/release", handler.Release)
	r.POST("/rollback", handler.Rollback)
	r.GET("/release/history", handler.ListReleaseHistory)

	return router
}

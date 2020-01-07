package main

import (
	"errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"

	_ "github.com/go-sql-driver/mysql"
	"XConf/config-srv/broadcast"
	"XConf/config-srv/broadcast/broker"
	"XConf/config-srv/broadcast/database"
	"XConf/config-srv/conf"
	"XConf/config-srv/dao"
	"XConf/config-srv/handler"
	protoConfig "XConf/proto/config"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"
)

var config conf.Config

func main() {
	log.Name("XConf")
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		// etcd的地址：http://127.0.0.1:2379
		op.Addrs = []string{"http://127.0.0.1:2379"}
	})
	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
		micro.Registry(reg),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_driver",
				Usage:  "database driver",
				EnvVar: "DATABASE_DRIVER",
				Value:  "mysql",
			},
			cli.StringFlag{
				Name:   "database_url",
				Usage:  "database url",
				EnvVar: "DATABASE_URL",
				Value:  "root:rootadmin@(127.0.0.1:3306)/xconf?charset=utf8&parseTime=true&loc=Local",
			},
			cli.StringFlag{
				Name:   "broadcast",
				Usage:  "broadcast (db/broker)",
				EnvVar: "BROADCAST",
				Value:  "db",
			}),
	)
	// 解析命令行标识参数可以使用service.Init
	service.Init(
		// 增加标识参数可以使用micro.Action选项
		// 也就是说如果有flag选项，那么解析flag选项需要使用micro.Action进行获取
		micro.Action(func(c *cli.Context) {
			config.DB.DriverName = c.String("database_driver")
			config.DB.URL = c.String("database_url")
			config.BroadcastType = c.String("broadcast")
			log.Infof("database_driver: %s , database_url: %s\n", config.DB.DriverName, config.DB.URL)
		}),
		// 服务start前的回调函数
		micro.BeforeStart(func() (err error) {
			if err = dao.Init(&config); err != nil {
				return
			}
			if err = dao.GetDao().Ping(); err != nil {
				return
			}

			var bc broadcast.Broadcast
			switch config.BroadcastType {
			case "db":
				bc, err = database.New()
				if err != nil {
					return err
				}
			case "broker":
				bc, err = broker.New(service)
				if err != nil {
					return err
				}
			default:
				return errors.New("broadcast： Invalid option")
			}
			broadcast.Init(bc)
			return
		}),
		// 服务stop前的回调函数
		micro.BeforeStop(func() error {
			return dao.GetDao().Disconnect()
		}),
	)

	if err := protoConfig.RegisterConfigHandler(service.Server(), new(handler.Config)); err != nil {
		panic(err)
	}

	if err := service.Run(); err != nil {
		panic(err)
	}
}

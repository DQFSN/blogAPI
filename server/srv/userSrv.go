package main

import (
	"github.com/DQFSN/blog/config"
	mpb "github.com/DQFSN/blog/proto/micro"
	"github.com/DQFSN/blog/server/microimp"
	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/v2/registry/consul"

	"log"
)

func main() {

	// 获取consul配置
	conf := config.Get().Consul
	address := conf.Host + ":" + conf.Port
	consulReg := consul.NewRegistry(
		registry.Addrs(address),
	)

	service := micro.NewService(
		micro.Name("user service"),
		micro.Registry(consulReg),
	)

	service.Init()

	err := mpb.RegisterAuthHandler(service.Server(), microimp.AuthHandler{})
	if err != nil {
		log.Fatal(err)
	}

	//运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

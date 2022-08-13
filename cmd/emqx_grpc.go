package main

import (
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"

	adapter "emqx-grpc/pkg/adapter/http"
	bridge "emqx-grpc/pkg/bridge/kafka"
	"emqx-grpc/pkg/entity"
)

func main() {

	configFile, err := ioutil.ReadFile("./configs/config.yaml")

	if err != nil {
		log.Panicln("读取配置文件失败", err)
	}
	var config *entity.Config
	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		log.Panicln("解析配置文件失败", err)
	}

	lis, err := net.Listen("tcp", config.Grpc.BindAddress)
	if err != nil {
		log.Fatalf("监听端口 %v 失败: %v", config.Grpc.BindAddress, err)
	}

	s := grpc.NewServer()

	adapterServer := adapter.RegisterGrpc(s, config.Emqx)
	bridgeServer := bridge.RegisterGrpc(s, config.Kafka)

	defer adapterServer.Close()
	defer bridgeServer.Close()

	log.Println("启动grpc服务 ", config.Grpc.BindAddress)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动grpc服务 %v 失败: %v", config.Grpc.BindAddress, err)
	}
}

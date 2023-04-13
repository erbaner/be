package main

import (
	"flag"
	"fmt"

	"github.com/erbaner/be/internal/rpc/msg"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImMessagePort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessagePrometheusPort[0], "msgPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start msg rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := msg.NewRpcChatServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}

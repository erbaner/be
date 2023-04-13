package main

import (
	"flag"
	"fmt"

	rpcAuth "github.com/erbaner/be/internal/rpc/auth"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAuthPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.AuthPrometheusPort[0], "authPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}

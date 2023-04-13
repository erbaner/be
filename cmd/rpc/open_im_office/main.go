package main

import (
	"flag"
	"fmt"

	rpc "github.com/erbaner/be/internal/rpc/office"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOfficePort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.OfficePrometheusPort[0], "officePrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start office rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpc.NewOfficeServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}

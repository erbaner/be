package main

import (
	"flag"
	"fmt"

	"github.com/erbaner/be/internal/rpc/organization"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOrganizationPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcOrganizationPort from cmd,default 11200 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.OrganizationPrometheusPort[0], "organizationPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start organization rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := organization.NewServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}

package main

import (
	"flag"
	"fmt"

	rpcMessageCMS "github.com/erbaner/be/internal/rpc/admin_cms"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAdminCmsPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.AdminCmsPrometheusPort[0], "adminCMSPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start cms rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcMessageCMS.NewAdminCMSServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}

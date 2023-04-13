package main

import (
	rpcCache "github.com/erbaner/be/internal/rpc/cache"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"

	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImCachePort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.CachePrometheusPort[0], "cachePrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start cache rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcCache.NewCacheServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}

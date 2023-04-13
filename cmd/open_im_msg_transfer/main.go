package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/erbaner/be/internal/msg_transfer/logic"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/log"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "MessageTransferPrometheusPort default listen port")
	flag.Parse()
	log.NewPrivateLog(constant.LogFileName)
	logic.Init()
	fmt.Println("start msg_transfer server ", ", OpenIM version: ", constant.CurrentVersion, "\n")
	logic.Run(*prometheusPort)
	wg.Wait()
}

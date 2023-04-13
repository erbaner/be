package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/erbaner/be/internal/cms_api"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/utils"

	"github.com/erbaner/be/pkg/common/config"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := cms_api.NewGinRouter()
	router.Use(utils.CorsHandler())
	defaultPorts := config.Config.CmsApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10006 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start cms api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion, "\n")
	router.Run(address)
}

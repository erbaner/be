/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	"fmt"

	pusher "github.com/erbaner/be/internal/push"
	fcm "github.com/erbaner/be/internal/push/fcm"
	"github.com/erbaner/be/internal/push/getui"
	jpush "github.com/erbaner/be/internal/push/jpush"
	"github.com/erbaner/be/internal/push/mobpush"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/kafka"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
	"github.com/erbaner/be/pkg/statistics"
)

var (
	rpcServer     RPCServer
	pushCh        PushConsumerHandler
	producer      *kafka.Producer
	offlinePusher pusher.OfflinePusher
	successCount  uint64
)

func Init(rpcPort int) {
	rpcServer.Init(rpcPort)
	pushCh.Init()

}
func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	statistics.NewStatistics(&successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	if *config.Config.Push.Getui.Enable {
		offlinePusher = getui.GetuiClient
	}
	if config.Config.Push.Jpns.Enable {
		offlinePusher = jpush.JPushClient
	}

	if config.Config.Push.Fcm.Enable {
		offlinePusher = fcm.NewFcm()
	}

	if config.Config.Push.Mob.Enable {
		offlinePusher = mobpush.MobPushClient
	}
}

func initPrometheus() {
	promePkg.NewMsgOfflinePushSuccessCounter()
	promePkg.NewMsgOfflinePushFailedCounter()
}

func Run(promethuesPort int) {
	go rpcServer.run()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
}

/*
** description("").
** copyright('Open_IM,www.Open_IM.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/13 10:33).
 */
package logic

import (
	"github.com/Shopify/sarama"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	kfk "github.com/erbaner/be/pkg/common/kafka"
	"github.com/erbaner/be/pkg/common/log"
	pbChat "github.com/erbaner/be/pkg/proto/msg"
	pbPush "github.com/erbaner/be/pkg/proto/push"
	"github.com/erbaner/be/pkg/utils"
	"github.com/golang/protobuf/proto"
)

type fcb func(msg []byte)

type PushConsumerHandler struct {
	msgHandle         map[string]fcb
	pushConsumerGroup *kfk.MConsumerGroup
}

func (ms *PushConsumerHandler) Init() {
	ms.msgHandle = make(map[string]fcb)
	ms.msgHandle[config.Config.Kafka.Ms2pschat.Topic] = ms.handleMs2PsChat
	ms.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ms2pschat.Topic}, config.Config.Kafka.Ms2pschat.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
}
func (ms *PushConsumerHandler) handleMs2PsChat(msg []byte) {
	log.NewDebug("", "msg come from kafka  And push!!!", "msg", string(msg))
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.Error("", "push Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	pbData := &pbPush.PushMsgReq{
		OperationID:  msgFromMQ.OperationID,
		MsgData:      msgFromMQ.MsgData,
		PushToUserID: msgFromMQ.PushToUserID,
	}
	sec := msgFromMQ.MsgData.SendTime / 1000
	nowSec := utils.GetCurrentTimestampBySecond()
	if nowSec-sec > 10 {
		return
	}
	switch msgFromMQ.MsgData.SessionType {
	case constant.SuperGroupChatType:
		MsgToSuperGroupUser(pbData)
	default:
		MsgToUser(pbData)
	}
	//Call push module to send message to the user
	//MsgToUser((*pbPush.PushMsgReq)(&msgFromMQ))
}
func (PushConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PushConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (ms *PushConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		ms.msgHandle[msg.Topic](msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

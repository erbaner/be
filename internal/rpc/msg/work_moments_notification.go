package msg

import (
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/log"
	pbOffice "github.com/erbaner/be/pkg/proto/office"
	sdk "github.com/erbaner/be/pkg/proto/sdk_ws"
	"github.com/erbaner/be/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func WorkMomentSendNotification(operationID, recvID string, notificationMsg *pbOffice.WorkMomentNotificationMsg) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), recvID, notificationMsg)
	WorkMomentNotification(operationID, recvID, recvID, notificationMsg)
}

func WorkMomentNotification(operationID, sendID, recvID string, m proto.Message) {
	var tips sdk.TipsComm
	var err error
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.UserMsgType,
		ContentType: constant.WorkMomentNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "proto.Marshal failed")
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), string(n.Content))
	Notification(n)
}

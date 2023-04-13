package msg

import (
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/log"

	//sdk "github.com/erbaner/be/pkg/proto/sdk_ws"
	"github.com/erbaner/be/pkg/utils"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"
)

func SuperGroupNotification(operationID, sendID, recvID string) {
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}

	log.NewInfo(operationID, utils.GetSelfFuncName(), string(n.Content))
	Notification(n)
}

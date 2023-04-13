/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/3/4 11:18).
 */
package im_mysql_msg_model

import (
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/db"
	"github.com/erbaner/be/pkg/common/log"
	pbMsg "github.com/erbaner/be/pkg/proto/msg"
	server_api_params "github.com/erbaner/be/pkg/proto/sdk_ws"
	"github.com/erbaner/be/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
)

func InsertMessageToChatLog(msg pbMsg.MsgDataToMQ) error {
	chatLog := new(db.ChatLog)
	copier.Copy(chatLog, msg.MsgData)
	switch msg.MsgData.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		chatLog.RecvID = msg.MsgData.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}
	if msg.MsgData.ContentType >= constant.NotificationBegin && msg.MsgData.ContentType <= constant.NotificationEnd {
		var tips server_api_params.TipsComm
		_ = proto.Unmarshal(msg.MsgData.Content, &tips)
		marshaler := jsonpb.Marshaler{
			OrigName:     true,
			EnumsAsInts:  false,
			EmitDefaults: false,
		}
		chatLog.Content, _ = marshaler.MarshalToString(&tips)

	} else {
		chatLog.Content = string(msg.MsgData.Content)
	}
	chatLog.CreateTime = utils.UnixMillSecondToTime(msg.MsgData.CreateTime)
	chatLog.SendTime = utils.UnixMillSecondToTime(msg.MsgData.SendTime)
	log.NewDebug("test", "this is ", chatLog)
	return db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").Create(chatLog).Error
}

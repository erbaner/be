package logic

import (
	"github.com/erbaner/be/pkg/common/db"
	"github.com/erbaner/be/pkg/common/log"
	pbMsg "github.com/erbaner/be/pkg/proto/msg"
	"github.com/erbaner/be/pkg/utils"
)

func saveUserChat(uid string, msg *pbMsg.MsgDataToMQ) error {
	time := utils.GetCurrentTimestampByMill()
	seq, err := db.DB.IncrUserSeq(uid)
	if err != nil {
		log.NewError(msg.OperationID, "data insert to redis err", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := pbMsg.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	log.NewInfo(msg.OperationID, "IncrUserSeq cost time", utils.GetCurrentTimestampByMill()-time)
	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
	//	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
}

func saveUserChatList(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (error, uint64) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", userID, len(msgList))
	//return db.DB.BatchInsertChat(userID, msgList, operationID)
	return db.DB.BatchInsertChat2Cache(userID, msgList, operationID)
}

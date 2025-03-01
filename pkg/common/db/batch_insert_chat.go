package db

import (
	"context"
	"errors"

	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/log"
	promePkg "github.com/erbaner/be/pkg/common/prometheus"
	pbMsg "github.com/erbaner/be/pkg/proto/msg"
	"github.com/erbaner/be/pkg/utils"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *DataBases) BatchInsertChat2DB(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string, currentMaxSeq uint64) error {
	newTime := getCurrentTimestampByMill()
	if len(msgList) > GetSingleGocMsgNum() {
		return errors.New("too large")
	}
	isInit := false
	var remain uint64
	blk0 := uint64(GetSingleGocMsgNum() - 1)
	//currentMaxSeq 4998
	if currentMaxSeq < uint64(GetSingleGocMsgNum()) {
		remain = blk0 - currentMaxSeq //1
	} else {
		excludeBlk0 := currentMaxSeq - blk0 //=1
		//(5000-1)%5000 == 4999
		remain = (uint64(GetSingleGocMsgNum()) - (excludeBlk0 % uint64(GetSingleGocMsgNum()))) % uint64(GetSingleGocMsgNum())
	}
	//remain=1
	insertCounter := uint64(0)
	msgListToMongo := make([]MsgInfo, 0)
	msgListToMongoNext := make([]MsgInfo, 0)
	seqUid := ""
	seqUidNext := ""
	log.Debug(operationID, "remain ", remain, "insertCounter ", insertCounter, "currentMaxSeq ", currentMaxSeq, userID, len(msgList))
	var err error
	for _, m := range msgList {
		log.Debug(operationID, "msg node ", m.String(), m.MsgData.ClientMsgID)
		currentMaxSeq++
		sMsg := MsgInfo{}
		sMsg.SendTime = m.MsgData.SendTime
		m.MsgData.Seq = uint32(currentMaxSeq)
		log.Debug(operationID, "mongo msg node ", m.String(), m.MsgData.ClientMsgID, "userID: ", userID, "seq: ", currentMaxSeq)
		if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
			return utils.Wrap(err, "")
		}
		if isInit {
			msgListToMongoNext = append(msgListToMongoNext, sMsg)
			seqUidNext = getSeqUid(userID, uint32(currentMaxSeq))
			log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
			continue
		}
		if insertCounter < remain {
			msgListToMongo = append(msgListToMongo, sMsg)
			insertCounter++
			seqUid = getSeqUid(userID, uint32(currentMaxSeq))
			log.Debug(operationID, "msgListToMongo ", seqUid, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain, "userID: ", userID)
		} else {
			msgListToMongoNext = append(msgListToMongoNext, sMsg)
			seqUidNext = getSeqUid(userID, uint32(currentMaxSeq))
			log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain, "userID: ", userID)
		}
	}

	ctx := context.Background()
	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)

	if seqUid != "" {
		filter := bson.M{"uid": seqUid}
		log.NewDebug(operationID, "filter ", seqUid, "list ", msgListToMongo, "userID: ", userID)
		err := c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgListToMongo}}}).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				filter := bson.M{"uid": seqUid}
				sChat := UserChat{}
				sChat.UID = seqUid
				sChat.Msg = msgListToMongo
				log.NewDebug(operationID, "filter ", seqUid, "list ", msgListToMongo)
				if _, err = c.InsertOne(ctx, &sChat); err != nil {
					promePkg.PromeInc(promePkg.MsgInsertMongoFailedCounter)
					log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
					return utils.Wrap(err, "")
				}
				promePkg.PromeInc(promePkg.MsgInsertMongoSuccessCounter)
			} else {
				promePkg.PromeInc(promePkg.MsgInsertMongoFailedCounter)
				log.Error(operationID, "FindOneAndUpdate failed ", err.Error(), filter)
				return utils.Wrap(err, "")
			}
		} else {
			promePkg.PromeInc(promePkg.MsgInsertMongoSuccessCounter)
		}
	}
	if seqUidNext != "" {
		filter := bson.M{"uid": seqUidNext}
		sChat := UserChat{}
		sChat.UID = seqUidNext
		sChat.Msg = msgListToMongoNext
		log.NewDebug(operationID, "filter ", seqUidNext, "list ", msgListToMongoNext, "userID: ", userID)
		if _, err = c.InsertOne(ctx, &sChat); err != nil {
			promePkg.PromeInc(promePkg.MsgInsertMongoFailedCounter)
			log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
			return utils.Wrap(err, "")
		}
		promePkg.PromeInc(promePkg.MsgInsertMongoSuccessCounter)
	}
	log.Debug(operationID, "batch mgo  cost time ", getCurrentTimestampByMill()-newTime, userID, len(msgList))
	return nil
}

func (d *DataBases) BatchInsertChat2Cache(insertID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (error, uint64) {
	newTime := getCurrentTimestampByMill()
	lenList := len(msgList)
	if lenList > GetSingleGocMsgNum() {
		return errors.New("too large"), 0
	}
	if lenList < 1 {
		return errors.New("too short as 0"), 0
	}
	// judge sessionType to get seq
	var currentMaxSeq uint64
	var err error
	if msgList[0].MsgData.SessionType == constant.SuperGroupChatType {
		currentMaxSeq, err = d.GetGroupMaxSeq(insertID)
		log.Debug(operationID, "constant.SuperGroupChatType  lastMaxSeq before add ", currentMaxSeq, "userID ", insertID, err)
	} else {
		currentMaxSeq, err = d.GetUserMaxSeq(insertID)
		log.Debug(operationID, "constant.SingleChatType  lastMaxSeq before add ", currentMaxSeq, "userID ", insertID, err)
	}
	if err != nil && err != go_redis.Nil {
		promePkg.PromeInc(promePkg.SeqGetFailedCounter)
		return utils.Wrap(err, ""), 0
	}
	promePkg.PromeInc(promePkg.SeqGetSuccessCounter)

	lastMaxSeq := currentMaxSeq
	for _, m := range msgList {

		currentMaxSeq++
		sMsg := MsgInfo{}
		sMsg.SendTime = m.MsgData.SendTime
		m.MsgData.Seq = uint32(currentMaxSeq)
		log.Debug(operationID, "cache msg node ", m.String(), m.MsgData.ClientMsgID, "userID: ", insertID, "seq: ", currentMaxSeq)
	}
	log.Debug(operationID, "SetMessageToCache ", insertID, len(msgList))
	err, failedNum := d.SetMessageToCache(msgList, insertID, operationID)
	if err != nil {
		promePkg.PromeAdd(promePkg.MsgInsertRedisFailedCounter, failedNum)
		log.Error(operationID, "setMessageToCache failed, continue ", err.Error(), len(msgList), insertID)
	} else {
		promePkg.PromeInc(promePkg.MsgInsertRedisSuccessCounter)
	}
	log.Debug(operationID, "batch to redis  cost time ", getCurrentTimestampByMill()-newTime, insertID, len(msgList))
	if msgList[0].MsgData.SessionType == constant.SuperGroupChatType {
		err = d.SetGroupMaxSeq(insertID, currentMaxSeq)
	} else {
		err = d.SetUserMaxSeq(insertID, currentMaxSeq)
	}
	if err != nil {
		promePkg.PromeInc(promePkg.SeqSetFailedCounter)
	} else {
		promePkg.PromeInc(promePkg.SeqSetSuccessCounter)
	}
	return utils.Wrap(err, ""), lastMaxSeq
}

//func (d *DataBases) BatchInsertChatBoth(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (error, uint64) {
//	err, lastMaxSeq := d.BatchInsertChat2Cache(userID, msgList, operationID)
//	if err != nil {
//		log.Error(operationID, "BatchInsertChat2Cache failed ", err.Error(), userID, len(msgList))
//		return err, 0
//	}
//	for {
//		if runtime.NumGoroutine() > 50000 {
//			log.NewWarn(operationID, "too many NumGoroutine ", runtime.NumGoroutine())
//			time.Sleep(10 * time.Millisecond)
//		} else {
//			break
//		}
//	}
//	return nil, lastMaxSeq
//}
//
//func (d *DataBases) BatchInsertChat(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) error {
//	newTime := getCurrentTimestampByMill()
//	if len(msgList) > GetSingleGocMsgNum() {
//		return errors.New("too large")
//	}
//	isInit := false
//	currentMaxSeq, err := d.GetUserMaxSeq(userID)
//	if err == nil {
//
//	} else if err == go_redis.Nil {
//		isInit = true
//		currentMaxSeq = 0
//	} else {
//		return utils.Wrap(err, "")
//	}
//	var remain uint64
//	//if currentMaxSeq < uint64(GetSingleGocMsgNum()) {
//	//	remain = uint64(GetSingleGocMsgNum()-1) - (currentMaxSeq % uint64(GetSingleGocMsgNum()))
//	//} else {
//	//	remain = uint64(GetSingleGocMsgNum()) - ((currentMaxSeq - (uint64(GetSingleGocMsgNum()) - 1)) % uint64(GetSingleGocMsgNum()))
//	//}
//
//	blk0 := uint64(GetSingleGocMsgNum() - 1)
//	if currentMaxSeq < uint64(GetSingleGocMsgNum()) {
//		remain = blk0 - currentMaxSeq
//	} else {
//		excludeBlk0 := currentMaxSeq - blk0
//		remain = (uint64(GetSingleGocMsgNum()) - (excludeBlk0 % uint64(GetSingleGocMsgNum()))) % uint64(GetSingleGocMsgNum())
//	}
//
//	insertCounter := uint64(0)
//	msgListToMongo := make([]MsgInfo, 0)
//	msgListToMongoNext := make([]MsgInfo, 0)
//	seqUid := ""
//	seqUidNext := ""
//	log.Debug(operationID, "remain ", remain, "insertCounter ", insertCounter, "currentMaxSeq ", currentMaxSeq, userID, len(msgList))
//	//4998 remain ==1
//	//4999
//	for _, m := range msgList {
//		log.Debug(operationID, "msg node ", m.String(), m.MsgData.ClientMsgID)
//		currentMaxSeq++
//		sMsg := MsgInfo{}
//		sMsg.SendTime = m.MsgData.SendTime
//		m.MsgData.Seq = uint32(currentMaxSeq)
//		if sMsg.Msg, err = proto.Marshal(m.MsgData); err != nil {
//			return utils.Wrap(err, "")
//		}
//		if isInit {
//			msgListToMongoNext = append(msgListToMongoNext, sMsg)
//			seqUidNext = getSeqUid(userID, uint32(currentMaxSeq))
//			log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
//			continue
//		}
//		if insertCounter < remain {
//			msgListToMongo = append(msgListToMongo, sMsg)
//			insertCounter++
//			seqUid = getSeqUid(userID, uint32(currentMaxSeq))
//			log.Debug(operationID, "msgListToMongo ", seqUid, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
//		} else {
//			msgListToMongoNext = append(msgListToMongoNext, sMsg)
//			seqUidNext = getSeqUid(userID, uint32(currentMaxSeq))
//			log.Debug(operationID, "msgListToMongoNext ", seqUidNext, m.MsgData.Seq, m.MsgData.ClientMsgID, insertCounter, remain)
//		}
//	}
//	//	ctx, _ := context.WithTimeout(context.Background(), time.Duration(config.Config.Mongo.DBTimeout)*time.Second)
//
//	ctx := context.Background()
//	c := d.mongoClient.Database(config.Config.Mongo.DBDatabase).Collection(cChat)
//
//	if seqUid != "" {
//		filter := bson.M{"uid": seqUid}
//		log.NewDebug(operationID, "filter ", seqUid, "list ", msgListToMongo)
//		err := c.FindOneAndUpdate(ctx, filter, bson.M{"$push": bson.M{"msg": bson.M{"$each": msgListToMongo}}}).Err()
//		if err != nil {
//			log.Error(operationID, "FindOneAndUpdate failed ", err.Error(), filter)
//			return utils.Wrap(err, "")
//		}
//	}
//	if seqUidNext != "" {
//		filter := bson.M{"uid": seqUidNext}
//		sChat := UserChat{}
//		sChat.UID = seqUidNext
//		sChat.Msg = msgListToMongoNext
//		log.NewDebug(operationID, "filter ", seqUidNext, "list ", msgListToMongoNext)
//		if _, err = c.InsertOne(ctx, &sChat); err != nil {
//			log.NewError(operationID, "InsertOne failed", filter, err.Error(), sChat)
//			return utils.Wrap(err, "")
//		}
//	}
//	log.NewWarn(operationID, "batch mgo  cost time ", getCurrentTimestampByMill()-newTime, userID, len(msgList))
//	return utils.Wrap(d.SetUserMaxSeq(userID, uint64(currentMaxSeq)), "")
//}

//func (d *DataBases)setMessageToCache(msgList []*pbMsg.MsgDataToMQ, uid string) (err error) {
//
//}

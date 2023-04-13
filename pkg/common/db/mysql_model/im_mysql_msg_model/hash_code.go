package im_mysql_msg_model

import (
	"hash/crc32"

	"github.com/erbaner/be/pkg/common/config"
)

func getHashMsgDBAddr(userID string) string {
	hCode := crc32.ChecksumIEEE([]byte(userID))
	return config.Config.Mysql.DBAddress[hCode%uint32(len(config.Config.Mysql.DBAddress))]
}

func getHashMsgTableIndex(userID string) int {
	hCode := crc32.ChecksumIEEE([]byte(userID))
	return int(hCode % uint32(config.Config.Mysql.DBMsgTableNum))
}

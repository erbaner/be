package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/erbaner/be/pkg/common/config"
)

func GetSign(paramsStr string) string {
	h := md5.New()
	io.WriteString(h, paramsStr)
	io.WriteString(h, config.Config.Push.Mob.AppSecret)
	fmt.Printf("%x", h.Sum(nil))

	return hex.EncodeToString(h.Sum(nil))
}

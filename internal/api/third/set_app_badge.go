package apiThird

import (
	"net/http"

	api "github.com/erbaner/be/pkg/base_info"
	"github.com/erbaner/be/pkg/common/db"
	"github.com/erbaner/be/pkg/common/log"
	"github.com/erbaner/be/pkg/common/token_verify"
	"github.com/erbaner/be/pkg/utils"
	"github.com/gin-gonic/gin"
)

func SetAppBadge(c *gin.Context) {
	var (
		req  api.SetAppBadgeReq
		resp api.SetAppBadgeResp
	)
	if err := c.Bind(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)

	var ok bool
	var errInfo, opUserID string
	ok, opUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	if !token_verify.CheckAccess(opUserID, req.FromUserID) {
		log.NewError(req.OperationID, "CheckAccess false ", opUserID, req.FromUserID)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "no permission"})
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req, opUserID)
	err := db.DB.SetUserBadgeUnreadCountSum(req.FromUserID, int(req.AppUnreadCount))
	if err != nil {
		errMsg := req.OperationID + " " + "SetUserBadgeUnreadCountSum failed " + err.Error() + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		resp.ErrCode = 500
		resp.ErrMsg = errMsg
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

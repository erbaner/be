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

func FcmUpdateToken(c *gin.Context) {
	var (
		req  api.FcmUpdateTokenReq
		resp api.FcmUpdateTokenResp
	)
	if err := c.Bind(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)

	ok, UserId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		resp.ErrCode = 500
		resp.ErrMsg = errMsg
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req, UserId)
	//逻辑处理开始
	err := db.DB.SetFcmToken(UserId, req.Platform, req.FcmToken, 0)
	if err != nil {
		errMsg := req.OperationID + " " + "SetFcmToken failed " + err.Error() + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		resp.ErrCode = 500
		resp.ErrMsg = errMsg
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	//逻辑处理完毕
	c.JSON(http.StatusOK, resp)
	return
}

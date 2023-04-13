package messageCMS

import (
	"context"
	"net/http"
	"strings"

	"github.com/erbaner/be/pkg/cms_api_struct"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/log"
	"github.com/erbaner/be/pkg/grpc-etcdv3/getcdv3"
	pbAdminCMS "github.com/erbaner/be/pkg/proto/admin_cms"
	pbCommon "github.com/erbaner/be/pkg/proto/sdk_ws"
	"github.com/erbaner/be/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetChatLogs(c *gin.Context) {
	var (
		req   cms_api_struct.GetChatLogsReq
		resp  cms_api_struct.GetChatLogsResp
		reqPb pbAdminCMS.GetChatLogsReq
	)
	if err := c.Bind(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}
	utils.CopyStructFields(&reqPb, &req)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdminCMS.NewAdminCMSClient(etcdConn)
	respPb, err := client.GetChatLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.ChatLogs {
		chatLog := cms_api_struct.ChatLog{}
		utils.CopyStructFields(&chatLog, v)
		resp.ChatLogs = append(resp.ChatLogs, &chatLog)
	}
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.ChatLogsNum = int(respPb.ChatLogsNum)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "resp", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}

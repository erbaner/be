package apiThird

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/client"
	"github.com/alibabacloud-go/tea/tea"
	api "github.com/erbaner/be/pkg/base_info"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/log"
	"github.com/erbaner/be/pkg/common/token_verify"
	"github.com/fatih/structs"

	//"github.com/fatih/structs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var stsClient *sts20150401.Client

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func getStsClient() *sts20150401.Client {
	if stsClient != nil {
		return stsClient
	}
	conf := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: tea.String(config.Config.Credential.Ali.AccessKeyID),
		// 您的AccessKey Secret
		AccessKeySecret: tea.String(config.Config.Credential.Ali.AccessKeySecret),
		// Endpoint
		Endpoint: tea.String(config.Config.Credential.Ali.StsEndpoint),
	}
	result, err := sts20150401.NewClient(conf)
	if err != nil {
		log.NewError("", "alists client初始化失败 ", err)
	}
	stsClient = result
	return stsClient
}

func AliOSSCredential(c *gin.Context) {
	req := api.OSSCredentialReq{}
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	var ok bool
	var userID string
	var errInfo string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "AliOSSCredential args ", userID)

	stsResp, err := getStsClient().AssumeRole(&sts20150401.AssumeRoleRequest{
		DurationSeconds: tea.Int64(config.Config.Credential.Ali.StsDurationSeconds),
		Policy:          nil,
		RoleArn:         tea.String(config.Config.Credential.Ali.OssRoleArn),
		RoleSessionName: tea.String(fmt.Sprintf("%s-%d", userID, time.Now().Unix())),
	})

	resp := api.OSSCredentialResp{}
	if err != nil {
		resp.ErrCode = constant.ErrTencentCredential.ErrCode
		resp.ErrMsg = err.Error()
	} else {
		resp = api.OSSCredentialResp{
			CommResp: api.CommResp{},
			OssData: api.OSSCredentialRespData{
				Endpoint:        config.Config.Credential.Ali.OssEndpoint,
				AccessKeyId:     *stsResp.Body.Credentials.AccessKeyId,
				AccessKeySecret: *stsResp.Body.Credentials.AccessKeySecret,
				Token:           *stsResp.Body.Credentials.SecurityToken,
				Bucket:          config.Config.Credential.Ali.Bucket,
				FinalHost:       config.Config.Credential.Ali.FinalHost,
			},
			Data: nil,
		}
	}

	resp.Data = structs.Map(&resp.OssData)
	log.NewInfo(req.OperationID, "AliOSSCredential return ", resp)

	c.JSON(http.StatusOK, resp)
}

package apiThird

import (
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
	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

func TencentCloudStorageCredential(c *gin.Context) {
	req := api.TencentCloudStorageCredentialReq{}
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
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "TencentCloudStorageCredential args ", userID)

	cli := sts.NewClient(
		config.Config.Credential.Tencent.SecretID,
		config.Config.Credential.Tencent.SecretKey,
		nil,
	)

	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          config.Config.Credential.Tencent.Region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:" + config.Config.Credential.Tencent.Region + ":uid/" + config.Config.Credential.Tencent.AppID + ":" + config.Config.Credential.Tencent.Bucket + "/*",
					},
				},
			},
		},
	}
	res, err := cli.GetCredential(opt)
	resp := api.TencentCloudStorageCredentialResp{}
	if err != nil {
		resp.ErrCode = constant.ErrTencentCredential.ErrCode
		resp.ErrMsg = err.Error()
	} else {
		resp.CosData.Bucket = config.Config.Credential.Tencent.Bucket
		resp.CosData.Region = config.Config.Credential.Tencent.Region
		resp.CosData.CredentialResult = res
	}

	resp.Data = structs.Map(&resp.CosData)
	log.NewInfo(req.OperationID, "TencentCloudStorageCredential return ", resp)

	c.JSON(http.StatusOK, resp)
}

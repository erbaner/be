package apiThird

import (
	"context"
	"net/http"

	api "github.com/erbaner/be/pkg/base_info"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/log"
	"github.com/erbaner/be/pkg/common/token_verify"
	"github.com/erbaner/be/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gin-gonic/gin"
)

func AwsStorageCredential(c *gin.Context) {
	var (
		req  api.AwsStorageCredentialReq
		resp api.AwsStorageCredentialResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	//原始帐号信息
	awsSourceConfig, err := awsConfig.LoadDefaultConfig(context.TODO(), awsConfig.WithRegion(config.Config.Credential.Aws.Region),
		awsConfig.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     config.Config.Credential.Aws.AccessKeyID,
				SecretAccessKey: config.Config.Credential.Aws.AccessKeySecret,
				Source:          "Open IM OSS",
			},
		}))
	if err != nil {
		errMsg := req.OperationID + " " + "Init AWS S3 Credential failed " + err.Error() + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	//帐号转化
	awsStsClient := sts.NewFromConfig(awsSourceConfig)
	StsRole, err := awsStsClient.AssumeRole(context.Background(), &sts.AssumeRoleInput{
		RoleArn:         aws.String(config.Config.Credential.Aws.RoleArn),
		DurationSeconds: aws.Int32(constant.AwsDurationTimes),
		RoleSessionName: aws.String(config.Config.Credential.Aws.RoleSessionName),
		ExternalId:      aws.String(config.Config.Credential.Aws.ExternalId),
	})
	if err != nil {
		errMsg := req.OperationID + " " + "AWS S3 AssumeRole failed " + err.Error() + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp.CosData.AccessKeyId = string(*StsRole.Credentials.AccessKeyId)
	resp.CosData.SecretAccessKey = string(*StsRole.Credentials.SecretAccessKey)
	resp.CosData.SessionToken = string(*StsRole.Credentials.SessionToken)
	resp.CosData.Bucket = config.Config.Credential.Aws.Bucket
	resp.CosData.RegionID = config.Config.Credential.Aws.Region
	resp.CosData.FinalHost = config.Config.Credential.Aws.FinalHost
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

package clientInit

import (
	"net/http"

	api "github.com/erbaner/be/pkg/base_info"
	imdb "github.com/erbaner/be/pkg/common/db/mysql_model/im_mysql_model"
	"github.com/erbaner/be/pkg/common/log"
	"github.com/erbaner/be/pkg/common/token_verify"
	"github.com/erbaner/be/pkg/utils"

	"github.com/gin-gonic/gin"
)

func SetClientInitConfig(c *gin.Context) {
	var req api.SetClientInitConfigReq
	var resp api.SetClientInitConfigResp
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	err, _ := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	m := make(map[string]interface{})
	if req.DiscoverPageURL != nil {
		m["discover_page_url"] = *req.DiscoverPageURL
	}
	if len(m) > 0 {
		err := imdb.SetClientInitConfig(m)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
			return
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetClientInitConfig(c *gin.Context) {
	var req api.GetClientInitConfigReq
	var resp api.GetClientInitConfigResp
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	err, _ := token_verify.ParseTokenGetUserID(c.Request.Header.Get("token"), req.OperationID)
	if err != nil {
		errMsg := "ParseTokenGetUserID failed " + err.Error() + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	config, err := imdb.GetClientInitConfig()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return

	}
	resp.Data.DiscoverPageURL = config.DiscoverPageURL
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp ", resp)
	c.JSON(http.StatusOK, resp)

}

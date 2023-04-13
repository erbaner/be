package friend

import (
	cbApi "github.com/erbaner/be/pkg/call_back_struct"
	"github.com/erbaner/be/pkg/common/config"
	"github.com/erbaner/be/pkg/common/constant"
	"github.com/erbaner/be/pkg/common/http"
	"github.com/erbaner/be/pkg/common/log"
	pbFriend "github.com/erbaner/be/pkg/proto/friend"

	//"github.com/erbaner/be/pkg/proto/msg"
	http2 "net/http"

	"github.com/erbaner/be/pkg/utils"
)

func callbackBeforeAddFriend(req *pbFriend.AddFriendReq) cbApi.CommonCallbackResp {
	callbackResp := cbApi.CommonCallbackResp{OperationID: req.CommID.OperationID}
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return callbackResp
	}
	log.NewDebug(req.CommID.OperationID, utils.GetSelfFuncName(), req.String())
	commonCallbackReq := &cbApi.CallbackBeforeAddFriendReq{
		CallbackCommand: constant.CallbackBeforeAddFriendCommand,
		FromUserID:      req.CommID.FromUserID,
		ToUserID:        req.CommID.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     req.CommID.OperationID,
	}
	resp := &cbApi.CallbackBeforeAddFriendResp{
		CommonCallbackResp: &callbackResp,
	}
	//utils.CopyStructFields(req, msg.MsgData)
	defer log.NewDebug(req.CommID.OperationID, utils.GetSelfFuncName(), commonCallbackReq, *resp)
	if err := http.CallBackPostReturn(config.Config.Callback.CallbackUrl, constant.CallbackBeforeAddFriendCommand, commonCallbackReq, resp, config.Config.Callback.CallbackBeforeAddFriend.CallbackTimeOut); err != nil {
		callbackResp.ErrCode = http2.StatusInternalServerError
		callbackResp.ErrMsg = err.Error()
		if !config.Config.Callback.CallbackBeforeAddFriend.CallbackFailedContinue {
			callbackResp.ActionCode = constant.ActionForbidden
			return callbackResp
		} else {
			callbackResp.ActionCode = constant.ActionAllow
			return callbackResp
		}
	}
	return callbackResp
}

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lqiz/mpai/app"
	oa "github.com/lqiz/mpai/pkg/official_account"
)

func WxMsg(routerGroup *gin.RouterGroup) {

	r := routerGroup.Group("/v1")
	acc := oa.NewExampleOfficialAccount(app.App.WC)

	r.Any("/serve", acc.Serve)

	//// 授权登录
	//r.GET("/join", acc.GetLoginAndShowPage)
	//r.GET("/login_redirect", acc.GetWxOauthRedirect)
	//// 新增永久图文素材
	//r.POST("/media/image/add", acc.UploadWxMedia)

	//获取ak
	r.GET("/oa/basic/get_access_token", acc.GetAccessToken)
	//获取微信callback IP
	r.GET("/oa/basic/get_callback_ip", acc.GetCallbackIP)
	//获取微信API接口 IP
	r.GET("/oa/basic/get_api_domain_ip", acc.GetAPIDomainIP)
	//清理接口调用次数
	r.GET("/api/v1/oa/basic/clear_quota", acc.ClearQuota)
}

type WxMsgEndpoint struct {
}

func NewWxMsgEndpoint() *WxMsgEndpoint {
	return &WxMsgEndpoint{}
}

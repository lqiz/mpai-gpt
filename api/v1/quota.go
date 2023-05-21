package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lqiz/mpai/api/request"
	"github.com/lqiz/mpai/pkg/quota"
	"github.com/lqiz/mpai/response"
	"net/http"
)

func Quota(routerGroup *gin.RouterGroup) {

	r := routerGroup.Group("/mp")
	endpoint := NewQuotaEndpoint()

	// 登录函数：登录接口调用该函数
	r.POST("/get-quota", endpoint.GetQuota)
	r.POST("/prompt", endpoint.Prompt)
}

type QuotaEndpoint struct {
	svc *quota.Quota
}

func NewQuotaEndpoint() *QuotaEndpoint {
	return &QuotaEndpoint{
		svc: quota.NewQuotaService(),
	}
}

// GetQuota 获取额度。
func (endpoint *QuotaEndpoint) GetQuota(c *gin.Context) {
	req := new(request.MpUserReq)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, response.ErrParams)
		return
	}

	quota := endpoint.svc.GetUserQuota(req.OpenId)

	c.JSON(http.StatusOK, response.NewOKResponse(map[string]int64{"quota": quota}))

}

// Prompt 开关
func (endpoint *QuotaEndpoint) Prompt(c *gin.Context) {
	req := new(request.PromptReq)
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusOK, response.ErrParams)
		return
	}
	result := endpoint.svc.OnOffSwitch(req.OnOff)

	c.JSON(http.StatusOK, response.NewOKResponse(map[string]int{"result": result}))
}

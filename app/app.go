package app

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/lqiz/mpai/model"
	"github.com/silenceper/wechat/v2"
	"gorm.io/gorm"
)

// App is the app singleton
var App *Application

type Application struct {
	Config        *model.Config
	OpenAiRouter  *gin.Engine
	RouterSwagger *gin.Engine
	DB            *gorm.DB
	Node          *snowflake.Node
	WC            *wechat.Wechat
	Client        *redis.Client
}

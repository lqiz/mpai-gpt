package main

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
	"github.com/lqiz/mpai/api"
	"github.com/lqiz/mpai/model"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lqiz/mpai/app"
	"github.com/lqiz/mpai/middleware"
)

func main() {
	app.App = &app.Application{}
	confPath := flag.String("c", "config.toml", "configure file")
	app.App.Config = model.GetConfig(confPath)

	flag.Parse()
	time.LoadLocation("Asia/Shanghai")

	// API 业务端口
	app.App.OpenAiRouter = gin.Default()
	app.App.OpenAiRouter.LoadHTMLFiles("template/index.html")
	app.App.OpenAiRouter.Static("/static", "template/static")
	app.App.OpenAiRouter.Use(middleware.WithTrace(), middleware.AccessLog())

	//获取wechat、redis 实例
	app.App.WC = InitWechat()
	app.App.Client = InitRedis()
	app.App.WC.SetCache(InitWXRedis())

	api.RouteMp(app.App.OpenAiRouter)
	log.SetLevel(log.DebugLevel)
	// 开始监听
	go func() {
		app.App.OpenAiRouter.Run(app.App.Config.Listen.Port)
	}()

	select {}
}

// InitWechat 获取wechat实例
// 在这里已经设置了全局cache，则在具体获取公众号/小程序等操作实例之后无需再设置，设置即覆盖
func InitWechat() *wechat.Wechat {
	cfg := app.App.Config
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        cfg.Redis.Host,
		Password:    cfg.Redis.Password,
		Database:    cfg.Redis.Database,
		MaxActive:   cfg.Redis.MaxActive,
		MaxIdle:     cfg.Redis.MaxIdle,
		IdleTimeout: cfg.Redis.IdleTimeout,
	}
	ctx := context.TODO()
	redisCache := cache.NewRedis(ctx, redisOpts)
	wc.SetCache(redisCache)
	return wc
}

// InitWXRedis 获取wechat实例
// 在这里已经设置了全局cache，则在具体获取公众号/小程序等操作实例之后无需再设置，设置即覆盖
func InitWXRedis() *cache.Redis {
	cfg := app.App.Config
	redisOpts := &cache.RedisOpts{
		Host:        cfg.Redis.Host,
		Password:    cfg.Redis.Password,
		Database:    cfg.Redis.Database,
		MaxActive:   cfg.Redis.MaxActive,
		MaxIdle:     cfg.Redis.MaxIdle,
		IdleTimeout: cfg.Redis.IdleTimeout,
	}
	ctx := context.TODO()
	redisCache := cache.NewRedis(ctx, redisOpts)

	return redisCache
}

func InitRedis() *redis.Client {
	cfg := app.App.Config
	options := &redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Database,
	}

	return redis.NewClient(options)
}

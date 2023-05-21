package official_account

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lqiz/mpai/app"
	"github.com/lqiz/mpai/dao"
	"github.com/lqiz/mpai/pkg"
	gpt "github.com/lqiz/mpai/pkg/gpt"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/singleflight"
	"strconv"
	"strings"
	"time"
)

// ExampleOfficialAccount 公众号操作样例
type ExampleOfficialAccount struct {
	wc              *wechat.Wechat
	officialAccount *officialaccount.OfficialAccount
	openAiClient    *gpt.OpenAIClient
	reqGroup        *singleflight.Group
	cacheMsg        *dao.MsgCache
}

// NewExampleOfficialAccount new
func NewExampleOfficialAccount(wc *wechat.Wechat) *ExampleOfficialAccount {
	//init config
	globalCfg := app.App.Config
	offCfg := &offConfig.Config{
		AppID:          globalCfg.AppID,
		AppSecret:      globalCfg.AppSecret,
		Token:          globalCfg.Token,
		EncodingAESKey: globalCfg.EncodingAESKey,
	}
	log.Debugf("offCfg=%+v", offCfg)
	officialAccount := wc.GetOfficialAccount(offCfg)
	return &ExampleOfficialAccount{
		wc:              wc,
		officialAccount: officialAccount,
		openAiClient:    gpt.NewOpenAIClient(),
		reqGroup:        &singleflight.Group{},
		cacheMsg:        dao.NewMsgCache(),
	}
}

// Serve 处理消息
func (ex *ExampleOfficialAccount) Serve(c *gin.Context) {
	ctx := c.Request.Context()
	// 传入request和responseWriter

	//u := ex.officialAccount.GetUser()
	server := ex.officialAccount.GetServer(c.Request, c.Writer)
	server.SkipValidate(true)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
		ctx2 := context.TODO()
		openID := msg.GetOpenID()

		// 【收到不支持的消息类型，暂无法显示】
		if strings.Contains(msg.Content, pkg.TipNoSupport) {
			log.Infof("%+v", msg)
			return decorateMsg(pkg.ErrNoSupport)
		}

		/*
			这里可以插入，判断是否有额度等商用代码，开源代码删除了，如果有需要私信沟通,
		*/

		count := ex.cacheMsg.GetCountCache(ctx2, openID)
		log.Infof("getCountCache1 %+v, time=%+v", count, time.Now().UnixNano())

		if count > 0 {
			replay, err := ex.cacheMsg.LoadMsgFromCache(ctx2, openID, count)
			if err != nil {
				return decorateMsg(pkg.ErrServer)
			}
			return decorateMsg(replay)
		}

		msgId := strconv.FormatInt(msg.MsgID, 10)
		result, err, _ := ex.reqGroup.Do(msgId, func() (interface{}, error) {
			return ex.callOpenAiWithTimeout(msg), nil
		})

		// 检测请求是否已经被中断或取消，非阻塞方式
		select {
		case <-ctx.Done():
			log.Infof("I AM CANCELED")
			// 请求已经被中断或取消，进行相应的清理工作
			return nil
		default:
			// 请求未被中断或取消，继续处理其他操作
		}

		log.Trace("callOpenAiWithTimeout result:=%+v, err=%+v, time=%+v", result, err, time.Now().UnixNano())

		if err != nil {
			return decorateMsg(pkg.ErrServer)
		}

		res := result.(bool)
		if res == false {
			return decorateMsg(pkg.TipWait)
		}

		// 立即回复
		count2 := ex.cacheMsg.GetCountCache(ctx, openID)
		log.Infof("getCountCache2 %+v time=%+v", count, time.Now().UnixNano())

		if count2 > 0 {
			replay, err := ex.cacheMsg.LoadMsgFromCache(ctx, openID, count2)
			if err != nil {
				return decorateMsg(pkg.ErrServer)
			}
			return decorateMsg(replay)
		}

		return nil

		//article1 := message.NewArticle("测试图文1", "图文描述", "", "")
		//articles := []*message.Article{article1}
		//news := message.NewNews(articles)
		//return &message.Reply{MsgType: message.MsgTypeNews, MsgData: news}

		//voice := message.NewVoice(mediaID)
		//return &message.Reply{MsgType: message.MsgTypeVoice, MsgData: voice}

		//
		//video := message.NewVideo(mediaID, "标题", "描述")
		//return &message.Reply{MsgType: message.MsgTypeVideo, MsgData: video}

		//music := message.NewMusic("标题", "描述", "音乐链接", "HQMusicUrl", "缩略图的媒体id")
		//return &message.Reply{MsgType: message.MsgTypeMusic, MsgData: music}

		//多客服消息转发
		//transferCustomer := message.NewTransferCustomer("")
		//return &message.Reply{MsgType: message.MsgTypeTransfer, MsgData: transferCustomer}
	})

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Errorf("Serve Error, err=%+v", err)
		return
	}
	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Errorf("Send Error, err=%+v", err)
		return
	}
}

func (ex *ExampleOfficialAccount) callOpenAiWithTimeout(msg *message.MixMessage) bool {
	result := make(chan bool, 1)
	ctxShort, _ := context.WithTimeout(context.Background(), 14*time.Second)
	ctxLong, cancelLong := context.WithTimeout(context.Background(), 30*time.Second)
	openID := msg.GetOpenID()

	go func() {
		select {
		case res := <-ex.openAiClient.CreateChatCompletion(ctxLong, msg.Content, openID):
			log.Infof("CreateChatCompletion finish res =%+v, err=%+v, time=%+v", res.GetText(), res.GetErr(), time.Now().UnixNano())
			content := res.GetText()
			if res.Err != nil {
				content = pkg.ErrServer
			} else {
				ex.cacheMsg.AddToList(ctxLong, openID, content)
			}
			ex.cacheMsg.CacheMsg(ctxLong, openID, content, time.Minute*2)

			result <- true
		case <-ctxLong.Done():
			fmt.Println("OPenAI timeout")
			ex.cacheMsg.CacheMsg(ctxLong, openID, pkg.ErrServer, time.Minute*2)
			cancelLong()
		}
	}()

	select {
	case <-result:
		fmt.Println("return true")
		return true
	case <-ctxShort.Done():
		fmt.Println("timeWait done")
		ex.cacheMsg.CacheMsg(ctxLong, openID, pkg.TipWait, time.Minute*2)
		return false
	}
}

func decorateMsg(text string) *message.Reply {
	log.Infof("decorateMsg + %+v", text)

	msg := message.NewText(text)
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: msg}
}

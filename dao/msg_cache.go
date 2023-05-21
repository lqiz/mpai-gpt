package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lqiz/mpai/app"
	"github.com/lqiz/mpai/pkg"
	log "github.com/sirupsen/logrus"
	"time"
)

const trunkSize = 500
const trunkTag = "_trunk"
const reReplayTip = "\n请回复【小鸮】继续查看"

type MsgCache struct {
	client *redis.Client
}

func NewMsgCache() *MsgCache {
	return &MsgCache{
		client: app.App.Client,
	}
}

func (mc *MsgCache) CacheMsg(ctx context.Context, key string, value string, expire time.Duration) {
	size := mc.chunkMessage(value)
	log.Infof("cacheMsg %+v, size = %+v, time=%+v", value, size, time.Now().UnixNano())

	cmd := app.App.Client.Set(ctx, key, value, expire)
	if cmd.Err() != nil {
		log.Infof("err = %+v", cmd.Err())
	}
	// 设置时间
	mc.SetCountCache(ctx, key, size, expire)
}

func (mc *MsgCache) SetCountCache(ctx context.Context, key string, count int, expire time.Duration) int {
	cmd := app.App.Client.Set(ctx, key+trunkTag, count, expire)
	err := cmd.Err()
	if err != nil {
		fmt.Println(err)
	}
	log.Infof("SetCountCache ret = %+v, time=%+v", count, time.Now().UnixNano())
	return count
}

func (mc *MsgCache) GetCountCache(ctx context.Context, key string) int {
	stringCmd := app.App.Client.Get(ctx, key+trunkTag)
	count, err := stringCmd.Int()
	if err != nil {
		fmt.Println(err)
	}
	log.Infof("GetCountCache ret = %+v, time=%+v", count, time.Now().UnixNano())
	return count
}

func (mc *MsgCache) decrCount(ctx context.Context, key string) {
	log.Infof("decrCount %+v, time=%+v", mc.GetCountCache(ctx, key), time.Now().UnixNano())
	cmd := app.App.Client.Decr(ctx, key+trunkTag)
	err := cmd.Err()
	if err != nil {
		log.Println(err)
	}
	return
}

func (mc *MsgCache) LoadMsgFromCache(ctx context.Context, key string, count int) (string, error) {
	stringCmd := app.App.Client.Get(ctx, key)
	result, err := stringCmd.Result()
	if err != nil {
		log.Infof("stringCmd.Result %+v", err)
	}

	// 如果还在思考，就直接返回
	if result == pkg.TipWait {
		return result, nil
	}

	// 如果count 大于 trunkSize，就分批次清理
	rr := []rune(result)
	len := len(rr)
	if (count-1)*trunkSize > len {
		log.Infof("count-1 * trunkSize  %+v", err)
		defer mc.delMsgCache(ctx, key)
		return "", errors.New("cache error")
	}

	size := mc.chunkMessage(result)
	start := size - count

	log.Infof("LoadMsgFromCache %+v, %+v, time=%+v", count, len, time.Now().UnixNano())

	max := (start + 1) * trunkSize
	if len < max {
		max = len
	}

	ret := rr[start*trunkSize : max]
	if count > 1 {
		ret = append(ret, []rune(reReplayTip)...)
	}

	log.Infof("RET %+v count=%+v", string(ret), count)
	// count 最后一个的时候，删除，其实不做也可以
	if count <= 1 {
		log.Infof("cont < = 1 = %+v", string(ret))
		defer mc.delMsgCache(ctx, key)
	}

	log.Infof("RET2 %+v", string(ret))

	// defer 后进先出
	mc.decrCount(ctx, key)

	return string(ret), nil
}

func (mc *MsgCache) delMsgCache(ctx context.Context, key string) {
	log.Infof("delMsgCache %+v", key)

	cmd := app.App.Client.Del(ctx, key)
	err := cmd.Err()
	if err != nil {
		log.Println(err)
	}

	tCmd := app.App.Client.Del(ctx, key+trunkTag)
	tErr := tCmd.Err()
	if tErr != nil {
		log.Println(err)
	}
	return
}

// 中文拆分
func (mc *MsgCache) chunkMessage(msg string) int {
	runes := []rune(msg)
	return len(runes)/trunkSize + 1
}

// 限制最近10条上下文，最大token 长度 1000
const contextMaxLength = 5
const maxToken = 1000
const contextCache = "context_"

// AddToList 插入 token 上下文 列表
func (mc *MsgCache) AddToList(ctx context.Context, key string, value string) error {
	client := app.App.Client
	key = contextCache + key

	// 获取列表长度
	lLenCmd := client.LLen(ctx, key)
	if lLenCmd.Err() != nil {
		return lLenCmd.Err()
	}
	listSize, err := lLenCmd.Result()
	if err != nil {
		return err
	}

	// 判断列表长度是否超过最大值，如果超过则移除最后一个元素，再添加新元素到列表头部
	if listSize >= contextMaxLength {
		if err := client.RPop(context.Background(), key).Err(); err != nil {
			return err
		}
		if err := client.LPush(context.Background(), key, value).Err(); err != nil {
			return err
		}
		return nil
	}

	client.Expire(ctx, key, time.Minute*5)

	// 如果列表长度未超过最大值，则直接将元素添加到列表头部
	if err := client.LPush(context.Background(), key, value).Err(); err != nil {
		return err
	}
	return nil
}

// GetListToken 获取 上下文 token  列表
func (mc *MsgCache) GetListToken(ctx context.Context, key string) ([][]rune, error) {
	client := app.App.Client
	key = contextCache + key

	listCmd := client.LRange(ctx, key, 0, -1)
	if listCmd.Err() != nil {
		return nil, listCmd.Err()
	}

	result, err := listCmd.Result()
	if err != nil {
		log.Error("GetListToken err=%+v", err)
		return nil, err
	}

	tokenLength := 0
	runeResult := make([][]rune, 0)
	for _, v := range result {
		rv := []rune(v)
		curLen := len(rv)
		remain := maxToken - tokenLength

		if remain >= curLen {
			runeResult = append(runeResult, rv)
			tokenLength += curLen
		} else {
			runeResult = append(runeResult, rv[0:remain])
			tokenLength = maxToken
		}
	}

	return runeResult, nil
}

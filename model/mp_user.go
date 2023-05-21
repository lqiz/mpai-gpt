package model

import (
	"fmt"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
	"github.com/silenceper/wechat/v2/officialaccount/user"
	"strings"
)

type MpUser struct {
	ID             int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Subscribe      int32  `gorm:"column:subscribe"`
	Openid         string `gorm:"column:openid"`
	Nickname       string `gorm:"column:nickname"`
	Sex            int32  `gorm:"column:sex"`
	City           string `gorm:"column:city"`
	Country        string `gorm:"column:country"`
	Province       string `gorm:"column:province"`
	Language       string `gorm:"column:language"`
	Headimgurl     string `gorm:"column:headimgurl"`
	SubscribeTime  int32  `gorm:"column:subscribe_time;default:0"`
	Unionid        string `gorm:"column:unionid"`
	Remark         string `gorm:"column:remark"`
	Groupid        int32  `gorm:"column:groupid"`
	TagidList      string `gorm:"column:tagid_list"`
	SubscribeScene string `gorm:"column:subscribe_scene"`
	QrScene        int    `gorm:"column:qr_scene"`
	QrSceneStr     string `gorm:"column:qr_scene_str"`
	Privilege      string `gorm:"column:privilege"`
}

func (m *MpUser) ToUser(info *user.Info) *MpUser {
	m.Groupid = info.GroupID
	m.Subscribe = info.Subscribe
	m.Openid = info.OpenID
	m.Nickname = info.Nickname
	m.Sex = info.Sex
	m.City = info.City
	m.Country = info.Country
	m.Province = info.Province
	m.Language = info.Language
	m.Headimgurl = info.Headimgurl
	m.SubscribeTime = info.SubscribeTime
	m.Unionid = info.UnionID
	m.Remark = info.Remark
	m.Groupid = info.GroupID
	m.TagidList = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(info.TagIDList)), ","), "[]")
	m.SubscribeScene = info.SubscribeScene
	m.QrScene = info.QrScene
	m.QrSceneStr = info.QrSceneStr
	return m
}

func (m *MpUser) ToUserOauth(info *oauth.UserInfo) *MpUser {
	m.Openid = info.OpenID
	m.Nickname = info.Nickname
	m.Sex = info.Sex
	m.Province = info.Province
	m.City = info.City
	m.Country = info.Country
	m.Headimgurl = info.HeadImgURL
	m.Privilege = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(info.Privilege)), ","), "[]")
	m.Unionid = info.Unionid
	return m
}

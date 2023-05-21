package model

type MpUserQuota struct {
	ID        int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	OpenID    string `gorm:"column:open_id;NOT NULL"`
	UnionID   string `gorm:"column:union_id"`
	Quota     int    `gorm:"column:quota;default:0;NOT NULL"` // 推广额度
	Gpt4Quota int    `gorm:"column:gpt4_quota;default:0;NOT NULL"`
}

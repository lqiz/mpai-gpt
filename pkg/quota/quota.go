package quota

type Quota struct {
	//mpUser:           dao.NewMpUser(),
	quotaSwitch bool
}

func NewQuotaService() *Quota {
	return &Quota{
		quotaSwitch: false,
	}
}

const defaultQuota = 1

// GetUserQuota 为了后面加每个用户限额方便，用户初始化时候赋予一定限额。
func (quotaSrv *Quota) GetUserQuota(openId string) int64 {
	if quotaSrv.quotaSwitch {
		return defaultQuota
	}

	return defaultQuota
}

func (quotaSrv *Quota) OnOffSwitch(status string) int {
	switch status {
	case "on":
		quotaSrv.quotaSwitch = true
		return 1
	case "off":
		quotaSrv.quotaSwitch = false
		return 2
	}

	return -1
}

package request

type LoginReq struct {
	Phone    string `json:"phone" form:"phone"`       // 手机号
	Email    string `json:"email" form:"email"`       // 邮箱
	Username string `json:"username" form:"username"` // 本来是使用格式验证，但是邮箱类型不止有一种
	Passwd   string `json:"passwd" form:"passwd"`
}

type RegisterReq struct {
	NickName string `json:"nick_name" form:"nick_name"`             // 昵称
	Phone    string `json:"phone" form:"phone" binding:"omitempty"` // 手机号
	Email    string `json:"email" form:"email" binding:"omitempty"` // 邮箱
	//Username string `json:"username" form:"username"`               // 本来是使用格式验证，但是邮箱类型不止有一种
	Passwd  string `json:"passwd" form:"passwd"`
	Captcha string `json:"captcha" form:"captcha"` // 验证码
}

type EmailCodeReq struct {
	Email string `json:"email" form:"email"` // 发送邮箱的验证码
}

type PhoneCodeReq struct {
	Phone string `json:"phone" form:"phone"` // 发送邮箱的验证码
}

type AuditReq struct {
	Email         string   `json:"email" form:"email"`                   // 待审核人员的邮箱
	Phone         string   `json:"phone" form:"phone"`                   // 待审核人员的手机号
	AuditApproved bool     `json:"audit_approved" form:"audit_approved"` // 审核通过
	SetAdmin      bool     `json:"set_admin" form:"set_admin"`           // 设置为管理员
	MineIdS       []string `json:"mine_id_str" form:"mine_id_str"`       // 绑定的矿井ID列表
}

type AuditListReq struct {
	Page    int64 `json:"page"`
	Size    int64 `json:"size"`
	IsAudit int   `json:"is_audit"`
}

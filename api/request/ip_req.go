package request

type APIKeysReq struct {
	IpAdd string `json:"ip_add" form:"ip_add"`
	UName string `json:"u_name" form:"u_name"`
}

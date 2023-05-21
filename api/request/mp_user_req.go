package request

type MpUserReq struct {
	UnionId string `json:"union_id" form:"union_id"` // union_id
	OpenId  string `json:"open_id" form:"open_id"`   // open_id
}

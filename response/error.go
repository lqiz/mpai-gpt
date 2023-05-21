package response

import "encoding/json"

var ErrorUnKnown = 4000

type Error struct {
	Status int    `json:"status"`
	Msg    string `json:"err_msg"`
}

func New(err int) *Error {
	return &Error{Status: err, Msg: "服务异常，请稍后重试"}
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Error) ErrNo() int {
	return e.Status
}

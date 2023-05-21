package response

import (
	"encoding/json"
	"time"
)

var (
	OKHeader  = NewHeader(0, "")
	ErrParams = NewErrorResponse(4000, "参数错误，请检查参数")
)

type Header struct {
	Code     int     `json:"c"` // status code
	Message  string  `json:"e"` // err msg
	Time     int64   `json:"s"` // timestamp
	Duration float32 `json:"t"` // latency
}

type Response struct {
	Header Header      `json:"h"`
	Data   interface{} `json:"c"`
}

type Body struct {
	Content interface{} `json:"c"`
}

func (e Header) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Header) Now() Header {
	e.Time = time.Now().Unix()
	return e
}

func NewHeader(code int, msg string) Header {
	return Header{
		Code:    code,
		Message: msg,
		Time:    time.Now().Unix(),
	}
}

func NewErrorResponse(code int, msg string) Response {
	return Response{
		Header: NewHeader(code, msg),
		Data:   struct{}{},
	}
}

func NewOKResponse(data interface{}) Response {
	return Response{
		Header: OKHeader.Now(),
		Data:   data,
	}
}

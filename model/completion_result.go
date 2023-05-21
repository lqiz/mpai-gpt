package model

import "github.com/lqiz/go-openai"

type Result struct {
	Val string
	Err error
}

type ResultGPT struct {
	Val *openai.ChatCompletionResponse
	Err error
}

func (r *ResultGPT) GetText() string {
	if r.Val == nil {
		return ""
	}
	return r.Val.Choices[0].Message.Content
}

func (r *ResultGPT) GetErr() error {
	return r.Err
}

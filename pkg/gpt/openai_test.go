package client

import (
	"context"
	"flag"
	"fmt"
	"github.com/lqiz/go-openai"
	"github.com/lqiz/mpai/app"
	"github.com/lqiz/mpai/model"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"testing"
)

func TestOpenAIClient_CreateChatCompletion(t *testing.T) {
	confPath := flag.String("c", "/Users/ruiyiluo/go/mpai/config.toml", "configure file")
	log.Infof("%+v", confPath)
	app.App = &app.Application{}
	app.App.Config = model.GetConfig(confPath)

	//client := NewOpenAIClient()

	config := openai.DefaultConfig("token")
	proxyUrl, err := url.Parse("http://165.154.134.147:4239")
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}

	config.BaseURL = "http://165.154.134.147:4239/v1"
	c := openai.NewClientWithConfig(config)

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "hello",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	log.Infof("%+v", resp)

}

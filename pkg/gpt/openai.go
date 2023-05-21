package client

import (
	"context"
	"fmt"
	"github.com/lqiz/go-openai"
	"github.com/lqiz/mpai/app"
	"github.com/lqiz/mpai/dao"
	"github.com/lqiz/mpai/model"
	"github.com/lqiz/mpai/pkg"
	log "github.com/sirupsen/logrus"
)

type OpenAIClient struct {
	client   *openai.Client
	cacheMsg *dao.MsgCache
}

func NewOpenAIClient() *OpenAIClient {
	config := openai.DefaultConfig("fake-token")
	config.BaseURL = app.App.Config.RemoteProxy.Url
	return &OpenAIClient{
		client:   openai.NewClientWithConfig(config),
		cacheMsg: dao.NewMsgCache(),
	}
}

func (ai *OpenAIClient) CreateChatCompletion(ctx context.Context, msg string, openID string) <-chan model.ResultGPT {
	ch := make(chan model.ResultGPT, 1)

	messages := ai.messageCompose(ctx, msg, openID)
	optionHeader := map[string]string{
		"X-UID": openID,
	}

	go func() {
		resp, err := ai.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				Messages:    messages,
				Temperature: 0.9,
				TopP:        1,
				MaxTokens:   1200,
				Stop:        []string{"\n\n\n"},
			},
			optionHeader,
		)

		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			ch <- model.ResultGPT{Err: err}
		}

		log.Infof("CreateChatCompletion return %+v", resp)
		ch <- model.ResultGPT{Val: &resp}
	}()

	return ch
}

func (ai *OpenAIClient) messageCompose(ctx context.Context, msg string, key string) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: pkg.SystemPrompt,
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		}}

	log.Infof(" messageCompose msg = %+v", string(msg))
	recentList, err := ai.cacheMsg.GetListToken(ctx, key)

	if err != nil {
		log.Errorf(" GetListToken err = %+v", err)
		return messages
	}

	for _, v := range recentList {
		msg := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: string(v),
		}
		log.Infof(" messageCompose = %+v", string(v))
		messages = append(messages, msg)
	}

	messages = append(messages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		})

	return messages
}

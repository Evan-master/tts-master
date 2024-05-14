package handler

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type answererOpenAI struct {
	base
}

func NewAnswererOpenAI() Handler {
	ao := &answererOpenAI{
		baseHander,
	}
	ao.name = "OpenAIAnswer"
	return ao
}

func (ao *answererOpenAI) Handle(ctx context.Context, field *Field) error {
	ao.onHandleStatusChange(OpenaiAnswerStart, field.OnHandleStatusChange)
	content := field.Content
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: content,
		},
	}
	response, err := ao.answer(ctx, messages)
	if err != nil {
		return fmt.Errorf("OpenAI ChatCompletion error: %v", err)
	}
	ao.onHandleStatusChange(OpenaiAnswerEnd, field.OnHandleStatusChange)
	field.Answer = response

	if ao.next != nil {
		return ao.next.Handle(ctx, field)
	}
	return nil
}

func (b *base) GetOpenAIToken() string {
	return b.cfg.OpenAIToken
}

func (a *answererOpenAI) answer(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	openaiToken := baseHander.GetOpenAIToken()
	a.openaiClient = openai.NewClient(openaiToken)
	resp, err := a.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo1106,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI ChatCompletion error: %v\n", err)
	}
	return resp.Choices[0].Message.Content, nil
}

func (a *answererOpenAI) streamAnswer(ctx context.Context, messages []openai.ChatCompletionMessage, onData func(fragment string)) error {
	seed := new(int)
	*seed = 666
	openaiToken := baseHander.GetOpenAIToken()
	a.openaiClient = openai.NewClient(openaiToken)
	stream, err := a.openaiClient.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-3.5-turbo-0125",
			Seed:  seed,
			Tools: []openai.Tool{
				{
					Type: openai.ToolTypeFunction,
					Function: openai.FunctionDefinition{
						Name:        "get_query_answer",
						Description: "Get the answer in a given query",
						Parameters: jsonschema.Definition{
							Type: jsonschema.Object,
							Properties: map[string]jsonschema.Definition{
								"query": {
									Type:        jsonschema.String,
									Description: "query that need to be answered, e.g. How to Use Function calling?",
								},
							},
							Required: []string{"query"},
						},
					},
				},
			},
			ToolChoice: "none",
			Messages:   messages,
		},
	)
	if err != nil {
		return errors.Join(err, fmt.Errorf("answer.Err.answerWithStream"))
	}
	defer stream.Close()
	for {
		recv, err1 := stream.Recv()
		if errors.Is(err1, io.EOF) {
			return nil
		}
		if err1 != nil {
			return err1
		}
		if len(recv.Choices) == 0 {
			return nil
		}
		onData(recv.Choices[0].Delta.Content)
	}
}

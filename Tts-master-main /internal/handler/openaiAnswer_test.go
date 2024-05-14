package handler

import (
	"context"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestAnswer(t *testing.T) {
	a := &answererOpenAI{}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: "Hello, who are you?",
		},
	}

	resp, err := a.answer(context.Background(), messages)
	t.Log("resp:", resp)
	if err != nil {
		t.Fatalf("answer() returned an error: %v", err)
	}

	if resp == "" {
		t.Errorf("answer() returned an empty response")
	}

}

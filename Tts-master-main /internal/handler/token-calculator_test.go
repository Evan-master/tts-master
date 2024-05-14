package handler

import (
	"context"
	"testing"
)

func TestNewTokenCalculator(t *testing.T) {
	calculator := NewTokenCalculator().(*TokenCalculator)
	field := &Field{
		Answer:          "asdsadasdsad",
		TranscribedText: "asdsadsadasd",
	}
	calculator.Handle(context.Background(), field)
	t.Log(field.GetPromptToken())
	t.Log(field.GetCompletionToken())
}

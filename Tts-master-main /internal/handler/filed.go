package handler

import (
	"github.com/sashabaranov/go-openai"
)

type tokenField struct {
	promptToken     int
	completionToken int
}

type Field struct {
	tokenField           *tokenField
	OnHandleStatusChange func(status HandleStatus)      `json:"-"`
	RequestID            string                         `json:"request_id"`
	Model                string                         `json:"model"`
	Content              string                         `json:"content"`
	Answer               string                         `json:"answer"`
	AnswerFragment       string                         `json:"answer_fragment"`
	Messages             []openai.ChatCompletionMessage `json:"messages"`
	TranscribedText      string                         `json:"transcribed_text"`
	Path                 string                         `json:"savePath"`
}

func (f *Field) GetPromptToken() int {
	if f.tokenField != nil {
		return f.tokenField.promptToken
	}
	return 0
}

func (f *Field) GetCompletionToken() int {
	if f.tokenField != nil {
		return f.tokenField.completionToken
	}
	return 0
}

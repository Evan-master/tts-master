package response

import (
	"wisebase/pkg/pager"
)

type List interface {
}

type WithListResp[T List] struct {
	Data []*T `json:"data"`
}

type WithListAndPageResp[T List] struct {
	Data []*T `json:"data"`
	pager.Pager
}

type DailyResp struct {
	Category  string `json:"category"`
	Value     int64  `json:"value"`
	CreatedAt int64  `json:"created_at"`
}

type TtsResp[T any] struct {
	Type   CompletionType `json:"type"`
	Status string         `json:"status,omitempty"`
	Field  T              `json:"field,omitempty"`
}

type ProcessingRequest struct {
	InputText       string
	GPTResponse     string
	TranscribedText string
}

type OpenAIAnswerResp struct {
	Answer string `json:"answer"`
}

type TokenResp struct {
	Prompt     int `json:"prompt"`
	Completion int `json:"complation"`
}

type Text2VoiceEnResp struct {
	TokenResult    string
	OpenAIAnswer   string
	Text2VoicePath string
}

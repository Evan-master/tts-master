package handler

import (
	"context"
)

type Handler interface {
	Name() string
	SetNext(handler Handler)
	Handle(ctx context.Context, field *Field) error
}

type HandleStatus string

func (h HandleStatus) String() string {
	return string(h)
}

const (
	ProgressStart       HandleStatus = "progress_start"
	Text2VoiceStart     HandleStatus = "Text2VoiceStart_start"
	Text2VoiceEnd       HandleStatus = "Text2VoiceStart_done"
	OpenaiAnswerStart   HandleStatus = "OpenaiAnswer_start"
	OpenaiAnswerEnd     HandleStatus = "OpenaiAnswer_done"
	Voice2TextStart     HandleStatus = "Voice2Text_start"
	Voice2TextEnd       HandleStatus = "Voice2Text_done"
	CalculateTokenStart HandleStatus = "calculate_token_start"
	CalculateTokenDone  HandleStatus = "calculate_token_done"
	ProgressDone        HandleStatus = "progress_done"
)

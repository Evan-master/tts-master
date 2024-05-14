package request

import (
	"github.com/sashabaranov/go-openai"

	"wisebase/pkg/random"
)

type TtsRequest struct {
	RequestID string                         `json:"-"`
	Model     string                         `json:"model"`
	Query     string                         `json:"query"`
	Messages  []openai.ChatCompletionMessage `json:"messages"`
}

func (s *TtsRequest) PreProcess() {
	s.genRequestID()
	s.removeSystemMessages()
}

func (s *TtsRequest) genRequestID() {
	s.RequestID = random.UUIDV4WithTimeStamp()
}

func (s *TtsRequest) removeSystemMessages() {
	if len(s.Messages) == 0 {
		return
	}
	msgs := make([]openai.ChatCompletionMessage, 0, len(s.Messages))
	for _, message := range s.Messages {
		if message.Role == openai.ChatMessageRoleSystem {
			continue
		}
		msgs = append(msgs, message)
	}
	s.Messages = msgs
}

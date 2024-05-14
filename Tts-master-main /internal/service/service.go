package service

import (
	"context"
	"wisebase/configs"
	"wisebase/internal/dto/request"
	"wisebase/internal/dto/response"
	"wisebase/internal/handler"
)

type TtsService struct {
	cfg             *configs.Config
	chain           handler.Handler
	fieldHandlerMap map[handler.HandleStatus]func(field *handler.Field) any
}

func NewVoiceService(cfg *configs.Config) *TtsService {
	handler.Init(&cfg.Handler)

	ts := &TtsService{
		cfg:             cfg,
		fieldHandlerMap: make(map[handler.HandleStatus]func(field *handler.Field) any),
	}
	ts.fieldHandlerMap[handler.OpenaiAnswerEnd] = ts.onOpenaiAnswerEnd
	ts.fieldHandlerMap[handler.Text2VoiceEnd] = ts.onText2VoiceEnd
	ts.fieldHandlerMap[handler.CalculateTokenDone] = ts.onCalculateTokenDone

	chain := handler.NewChain()
	tokenCalculator := handler.NewTokenCalculator()
	openAIAnswerer := handler.NewAnswererOpenAI()
	text2voice := handler.NewText2Voice()

	chain.SetNext(tokenCalculator)
	tokenCalculator.SetNext(openAIAnswerer)
	openAIAnswerer.SetNext(text2voice)
	ts.chain = chain

	return ts
}

func (s *TtsService) TextVoice(ctx context.Context, req *request.TtsRequest) (*response.Text2VoiceEnResp, error) {

	field := &handler.Field{
		RequestID: req.RequestID,
		Model:     req.Model,
		Content:   req.Query, // 将 query 设置到 Field 中
		Messages:  req.Messages,
	}

	err := s.chain.Handle(ctx, field)
	if err != nil {
		return nil, err
	}

	// 从 field 中读取结果
	tokenResult := field.Answer
	openAIAnswer := field.Answer
	text2VoicePath := field.Path

	rv := &response.Text2VoiceEnResp{
		TokenResult:    tokenResult,
		OpenAIAnswer:   openAIAnswer,
		Text2VoicePath: text2VoicePath,
	}
	return rv, nil
}

func (s *TtsService) onHandleStatusChange(handleStatus handler.HandleStatus, field *handler.Field) any {
	rv := response.TtsResp[any]{
		Type:   response.CompletionTypeText,
		Status: handleStatus.String(),
		Field:  nil,
	}
	f, ok := s.fieldHandlerMap[handleStatus]
	if ok {
		rv.Field = f(field)
	}
	return rv
}

func (s *TtsService) onText2VoiceEnd(field *handler.Field) any {
	return &response.Text2VoiceEnResp{
		Text2VoicePath: field.Path,
	}
}

func (s *TtsService) onOpenaiAnswerEnd(field *handler.Field) any {
	return &response.OpenAIAnswerResp{
		Answer: field.Answer,
	}
}

func (s *TtsService) onCalculateTokenDone(field *handler.Field) any {
	return &response.TokenResp{
		Prompt:     field.GetPromptToken(),
		Completion: field.GetCompletionToken(),
	}
}

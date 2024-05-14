package handler

import (
	"context"
)

type TokenCalculator struct {
	base
}

func NewTokenCalculator() Handler {
	t := &TokenCalculator{
		baseHander,
	}
	t.name = "TokenCalculator"
	return t
}

func (h *TokenCalculator) Handle(ctx context.Context, field *Field) error {
	defer func() {
		h.onHandleStatusChange(CalculateTokenStart, field.OnHandleStatusChange)
		h.handle(field)
		h.onHandleStatusChange(CalculateTokenDone, field.OnHandleStatusChange)
	}()

	if h.next == nil {
		return nil
	}
	return h.next.Handle(ctx, field)

}

func (h *TokenCalculator) handle(field *Field) {
	_, promptToken := h.encodeToken(answerPrompt + field.Answer + field.TranscribedText)
	_, completionToken := h.encodeToken(field.Answer)
	field.tokenField = &tokenField{
		promptToken:     promptToken,
		completionToken: completionToken,
	}
}

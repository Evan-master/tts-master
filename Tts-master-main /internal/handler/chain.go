package handler

import (
	"context"
)

type Chain struct {
	base
}

func (h *Chain) Handle(ctx context.Context, field *Field) error {
	h.onHandleStatusChange(ProgressStart, field.OnHandleStatusChange)
	defer h.onHandleStatusChange(ProgressDone, field.OnHandleStatusChange)
	if h.next == nil {
		return nil
	}
	return h.next.Handle(ctx, field)
}

func NewChain() Handler {
	cd := &Chain{base: baseHander}
	cd.name = "TTsMaster"
	return cd
}

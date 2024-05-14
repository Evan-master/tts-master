package handler

import (
	"context"
	"testing"
)

func TestHandle(t *testing.T) {
	a := NewVoice2Text()
	b := NewAnswererOpenAI()
	c := NewText2Voice()
	d := NewTokenCalculator()
	a.SetNext(b)
	b.SetNext(c)
	c.SetNext(d)

	field := &Field{}

	ctx := context.Background()

	err := a.Handle(ctx, field)
	t.Log(field)
	if err != nil {
		t.Errorf("Handle chain failed: %v", err)
	}
}

package handler

import (
	"testing"
)

func TestVoice2Text(t *testing.T) {
	response, err := Voice2Text()
	t.Log("Voice2Text() returned:", response)
	if err != nil {
		t.Errorf("Voice2Text() returned an error: %v", err)
	}

	expectedResponse := "Today is a wonderful day to build something people love."

	if response != expectedResponse {
		t.Errorf("Voice2Text() returned unexpected response: %s", response)
	}
}

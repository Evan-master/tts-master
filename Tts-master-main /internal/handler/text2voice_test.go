package handler

import (
	"os"
	"testing"
)

func TestText2voice(t *testing.T) {
	req := &SpeechRequest{
		Model: "tts-1",
		Voice: "alloy",
		Input: "你好，你可以教我学英语吗!",
	}

	_, err := Text2voice(req)
	if err != nil {
		t.Errorf("Text2voice() error = %v", err)
	}
	const audioPath = "save/answer.mp3"
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		t.Errorf("文件未成功保存：%s", audioPath)
	} else {
		t.Logf("文件已成功保存：%s", audioPath)
	}
}

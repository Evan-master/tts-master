package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type SpeechRequest struct {
	Model string `json:"model"`
	Voice string `json:"voice"`
	Input string `json:"input"`
}

type toVoice struct {
	base
}

func (b *base) GetOpenAIProxyHost1() string {
	return b.cfg.OpenAIProxyHost1
}

func (b *base) GetOpenAIToken1() string {
	return b.cfg.OpenAIToken
}

func NewText2Voice() Handler {
	tv := &toVoice{
		baseHander,
	}
	return tv
}

func (tv *toVoice) Handle(ctx context.Context, field *Field) error {
	tv.onHandleStatusChange(Text2VoiceStart, field.OnHandleStatusChange)
	speechReq := &SpeechRequest{
		Model: "tts-1-hd",
		Voice: "alloy",
		Input: field.Answer,
	}
	_, err := Text2voice(speechReq)
	if err != nil {
		return fmt.Errorf("text to voice conversion failed: %v", err)
	}
	tv.onHandleStatusChange(Text2VoiceEnd, field.OnHandleStatusChange)
	if tv.next != nil {
		return tv.next.Handle(ctx, field)
	}
	return nil
}

func Text2voice(s *SpeechRequest) (string, error) {
	openaiURL := baseHander.GetOpenAIProxyHost1()
	openaiToken := baseHander.GetOpenAIToken1()

	requestBody, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}
	log.Println("requestBody:", string(requestBody))

	req, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+openaiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK response status: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	answerDir := "save"
	if _, err := os.Stat(answerDir); os.IsNotExist(err) {
		if mkErr := os.Mkdir(answerDir, 0755); mkErr != nil {
			return "", fmt.Errorf("error creating directory: %v", mkErr)
		}
	}

	answerFilePath := filepath.Join(answerDir, "answer.mp3")
	err = ioutil.WriteFile(answerFilePath, body, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	return answerFilePath, nil
}

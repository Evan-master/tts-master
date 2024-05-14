package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type toText struct {
	base
}

// GetOpenAIProxyHost 提供对 OpenAIProxyHost 的访问
func (b *base) GetOpenAIProxyHost2() string {
	return b.cfg.OpenAIProxyHost2
}

// GetOpenAIToken 提供对 OpenAIToken 的访问
func (b *base) GetOpenAIToken2() string {
	return b.cfg.OpenAIToken
}

func NewVoice2Text() Handler {
	tt := &toText{
		baseHander,
	}
	tt.name = "Voice2Text"
	return tt
}

func (tt *toText) Handle(ctx context.Context, field *Field) error {
	tt.onHandleStatusChange(Voice2TextStart, field.OnHandleStatusChange)
	transText, err := Voice2Text()
	field.TranscribedText = transText
	if err != nil {
		return fmt.Errorf("voice to text conversion failed: %v", err)
	}
	tt.onHandleStatusChange(Voice2TextEnd, field.OnHandleStatusChange)
	// 如果存在下一个处理器，继续处理链
	if tt.next != nil {
		return tt.next.Handle(ctx, field)
	}
	return nil
}

func Voice2Text() (string, error) {
	openaiURL := baseHander.GetOpenAIProxyHost2()
	openaiToken := baseHander.GetOpenAIToken2()

	const questionFilePath = "save/question.mp3"

	file, err := os.Open(questionFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening audio file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("error copying file content: %v", err)
	}

	_ = writer.WriteField("model", "whisper-1")

	err = writer.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", openaiURL, body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+openaiToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK response status: %s", resp.Status)
	}

	var result struct {
		Text string `json:"text"`
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return "", fmt.Errorf("error parsing response JSON: %v", err)
	}

	return result.Text, nil
}

package handler

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	OpenAIToken      string `yaml:"OpenAIToken"`
	OpenAIProxyHost1 string `yaml:"OpenAIProxyHost1"`
	OpenAIProxyHost2 string `yaml:"OpenAIProxyHost2"`
	Timeout          int    `yaml:"Timeout"`
}

var baseHander = newBase(&Config{
	OpenAIProxyHost1: "https://api.openai.com/v1/audio/speech",
	OpenAIProxyHost2: "https://api.openai.com/v1/audio/transcriptions",
	OpenAIToken:      "",
	Timeout:          60,
})

type base struct {
	name         string
	cfg          Config
	httpClient   *resty.Client
	openaiClient *openai.Client
	tiktoken     *tiktoken.Tiktoken

	next Handler
}

func Init(cfg *Config) {
	baseHander = newBase(cfg)
}

func newBase(cfg *Config) base {
	if cfg == nil {
		panic("please input config")
	}
	if cfg.OpenAIProxyHost1 == "" {
		panic("please input OpenAIProxyHost1 config")
	}
	_, err := url.Parse(cfg.OpenAIProxyHost1)
	if cfg.OpenAIProxyHost2 == "" {
		panic("please input OpenAIProxyHost2 config")
	}
	_, err = url.Parse(cfg.OpenAIProxyHost2)

	if err != nil {
		panic(err)
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60
	}
	tiktokenClient, err := tiktoken.EncodingForModel("gpt-3.5-turbo")
	if err != nil {
		panic(err)
	}

	timeout := time.Duration(cfg.Timeout) * time.Second

	httpClient := resty.New().SetTimeout(timeout)

	openaiCfg := openai.DefaultConfig(cfg.OpenAIToken)
	openaiCfg.HTTPClient.Timeout = timeout
	openaiClient := openai.NewClientWithConfig(openaiCfg)

	return base{
		name:         "base",
		cfg:          *cfg,
		httpClient:   httpClient,
		openaiClient: openaiClient,
		tiktoken:     tiktokenClient,
	}
}

func (b *base) Name() string {
	return b.name
}

func (b *base) SetNext(handler Handler) {
	b.next = handler
}

func (b *base) handleError(err error) error {
	return errors.Join(err, fmt.Errorf("Handler.%s.Err", b.Name()))
}

func (b *base) onHandleStatusChange(status HandleStatus, handler func(HandleStatus)) {
	if handler != nil {
		handler(status)
	}
}

func (b *base) encodeToken(content string) ([]int, int) {
	encode := b.tiktoken.Encode(content, nil, nil)
	return encode, len(encode)
}

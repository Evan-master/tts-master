// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"wisebase/configs"
	"wisebase/internal/cmd"
	"wisebase/internal/route/v1"
	"wisebase/internal/service"
)

// Injectors from wire.go:

func build() (*cmd.App, func(), error) {
	config, err := configs.InitConfig()
	if err != nil {
		return nil, nil, err
	}
	engine := cmd.NewGIN(config)
	ttsService := service.NewVoiceService(config)
	option := &v1.Option{
		TtsService: ttsService,
	}
	v := v1.NewRoutes(option)
	cmdOption := &cmd.Option{
		Config:  config,
		Handler: engine,
		Routers: v,
	}
	app := cmd.NewApp(cmdOption)
	return app, func() {
	}, nil
}
//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	v1 "wisebase/internal/route/v1"

	"wisebase/configs"
	"wisebase/internal/cmd"
	"wisebase/internal/service"
)

func build() (*cmd.App, func(), error) {
	panic(wire.Build(
		configs.InitConfig,
		service.ProviderSet,
		v1.ProviderSet,
		cmd.ProviderSet,
	))
}

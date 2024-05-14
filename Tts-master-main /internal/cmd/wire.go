package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"wisebase/configs"
	v1 "wisebase/internal/route/v1"
)

type Option struct {
	Config  *configs.Config
	Handler *gin.Engine
	Routers []v1.Router
}

var ProviderSet = wire.NewSet(
	wire.Struct(new(Option), "*"),
	NewGIN,
	NewApp,
)

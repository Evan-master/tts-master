package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
	"wisebase/internal/service"

	"wisebase/configs"
)

type WireOption struct {
	Log     *zap.SugaredLogger
	Conf    *configs.Config
	Handler *gin.Engine

	Routers []Router
}

var ProviderSet = wire.NewSet(
	wire.Struct(new(Option), "*"),
	NewRoutes,
)

type Option struct {
	TtsService *service.TtsService
}

package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wisebase/configs"
	"wisebase/internal/dto/request"
	"wisebase/internal/route/common"
	"wisebase/internal/route/middleware"
	"wisebase/internal/service"
)

func NewGinEngine(conf *configs.Config) *gin.Engine {
	if conf.IsReleaseMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	f, _ := os.OpenFile("./log/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.Cors(true),
	)
	return r
}

type HttpEngine struct {
	log     *zap.SugaredLogger
	conf    *configs.Config
	handler *gin.Engine
	route   *Route // 添加 Route 实例
	routers []Router
}

func NewHttpEngine(opt *WireOption) *HttpEngine {
	return &HttpEngine{
		log:     opt.Log,
		conf:    opt.Conf,
		handler: opt.Handler,
		route:   &Route{srv: &service.TtsService{}},
		routers: opt.Routers,
	}
}

func NewRoutes(opt *Option) []Router {
	rv := make([]Router, 0)
	rv = append(rv, NewRoute(opt))
	return rv
}

func NewRoute(opt *Option) *Route {
	return &Route{
		srv: opt.TtsService,
	}
}

type Router interface {
	RegisterRoute(r *gin.RouterGroup)
}

type Route struct {
	srv *service.TtsService
}

func (r *Route) RegisterRoute(router *gin.RouterGroup) {
	router.POST("/tts", r.text2VoiceHandler)
}

func (r *Route) text2VoiceHandler(ctx *gin.Context) {
	var req request.TtsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 打印日志：绑定JSON失败
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}
	// 打印日志：请求数据绑定成功，准备预处理
	log.Printf("Request bound to struct: %+v", req)

	req.PreProcess()

	log.Printf("Request preprocessed: %+v", req)
	if r.srv == nil {
		// 打印日志：服务实例是 nil
		log.Println("Service instance is nil")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	// 打印日志：调用服务的 TextVoice 方法
	log.Printf("Calling TextVoice method with request: %+v", req)
	data, err := r.srv.TextVoice(ctx.Request.Context(), &req)
	//data, err := r.srv.TextVoice(ctx.Request.Context(), &req)
	if err != nil {
		// 打印日志：TextVoice 调用出错
		log.Printf("Error calling TextVoice: %v", err)
	}

	// 使用 WrapResp 处理结果和错误
	respFunc := common.WrapResp(ctx)
	respFunc(data, err) // 将 data 和 err 传递给返回的函数

}

func (h *HttpEngine) registerRoute() {
	// 创建 API v1 的路由组
	v1Group := h.handler.Group("/api/v1")

	// 注册 Route 实例的路由
	h.route.RegisterRoute(v1Group)

	// 注册其他路由器的路由
	for _, router := range h.routers {
		router.RegisterRoute(v1Group)
	}
}

func (h *HttpEngine) Run() error {
	common.SetRespLog(h.log)
	h.registerRoute()

	srv := &http.Server{
		Addr:    h.conf.App.Addr,
		Handler: h.handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		h.log.Fatal("Server Shutdown:", err)
		return err
	}
	h.log.Infof("server exiting")
	return nil
}

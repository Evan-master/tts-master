package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"wisebase/configs"
	"wisebase/internal/route/middleware"
	v1 "wisebase/internal/route/v1"
	"wisebase/pkg/xslog"
)

func NewGIN(conf *configs.Config) *gin.Engine {
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

type App struct {
	conf    *configs.Config
	handler *gin.Engine
	routers []v1.Router
}

func NewApp(opt *Option) *App {
	app := &App{conf: opt.Config, handler: opt.Handler, routers: opt.Routers}
	app.setDefaultSlog()
	return app
}

func (a *App) setDefaultSlog() {
	var extraWriters []xslog.ExtraWriter
	if a.conf.IsDebugMode() {
		a.conf.Log.Level = slog.LevelDebug
		extraWriters = append(extraWriters, xslog.ExtraWriter{
			Writer: os.Stdout,
			Level:  slog.LevelDebug,
		})
	}
	a.conf.Log.ExtraWriters = extraWriters

	fileLogger := xslog.NewFileSlog(&a.conf.Log)
	slog.SetDefault(fileLogger)
}

func (a *App) Run() error {
	a.registerRoute()

	srv := &http.Server{
		Addr:    a.conf.App.Addr,
		Handler: a.handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) registerRoute() {
	r := a.handler.Group("/api/v1")
	for _, router := range a.routers {
		router.RegisterRoute(r)
	}
}

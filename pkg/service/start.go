package service

import (
	glog "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/xuxant/Darwin_API/pkg/config"
	"github.com/xuxant/Darwin_API/pkg/logs"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"net/http"
)

func InitGinEngine(cfg config.ServiceConfig, logger zerolog.Logger) *gin.Engine {
	middlewareList := []gin.HandlerFunc{
		gin.Recovery(),
		gintrace.Middleware(cfg.ServiceID),
		logs.ZerologMiddleware(logger),
		panicHandlerMiddleWare(),
	}

	if !cfg.IsLocalProfile() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.TestMode)
		middlewareList = append(middlewareList,
			glog.SetLogger(glog.WithLogger(func(_ *gin.Context, l zerolog.Logger) zerolog.Logger {
				return l.Output(gin.DefaultWriter).With().Logger()
			}),
			))
	}

	g := gin.New()
	middlewares := CombineMiddlewares(middlewareList)
	g.Use(middlewares...)
	return g
}

func CombineMiddlewares(middlewareList ...[]gin.HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, 0)

	for _, handlers := range middlewareList {
		funcs = append(handlers)
	}

	return funcs
}

func AddHealthCheckRoute(group *gin.RouterGroup) gin.IRoutes {
	return group.GET("/healthcheck/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pinged")
	})
}

func panicHandlerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Status(500)
				c.Writer.WriteHeaderNow()
				panic(err)
			}
		}()
		c.Next()
	}
}

func Initialize(cfg config.ServiceConfig, logger zerolog.Logger) Server {
	engine := NewGin(cfg, logger)
	server := NewHttpServer(cfg, engine, logger)
	return server
}

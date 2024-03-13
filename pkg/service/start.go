package service

import (
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	glog "github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/xuxant/Darwin_API/pkg/config"
	"github.com/xuxant/Darwin_API/pkg/logs"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"net/http"
)

func InitGinEngine(cfg config.ServiceConfig, logger zerolog.Logger) (*gin.Engine, []gin.HandlerFunc) {
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

	if cfg.ClerkAuthentication {
		clerk.SetKey(cfg.ClerkSecretKey)
		clientConfig := &clerk.ClientConfig{}
		clientConfig.Key = &cfg.ClerkSecretKey
		jwksClient := jwks.NewClient(clientConfig)
		middlewareList = append(middlewareList, authHandlerMiddleware(jwksClient))
	}

	g := gin.New()
	middlewares := CombineMiddlewares(middlewareList)
	return g, middlewares
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

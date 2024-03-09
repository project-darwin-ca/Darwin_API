package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/xuxant/Darwin_API/pkg/config"
	"github.com/xuxant/Darwin_API/pkg/db"
	"github.com/xuxant/Darwin_API/pkg/handlers"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	srv             *http.Server
	logger          zerolog.Logger
	shutdownTimeOut int
}

func NewHttpServer(cfg config.ServiceConfig, ginEngine *gin.Engine, logger zerolog.Logger) Server {
	return Server{
		logger:          logger,
		shutdownTimeOut: cfg.ShutdownTimeout,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: ginEngine,
		},
	}
}

func NewGin(cfg config.ServiceConfig, logger zerolog.Logger) *gin.Engine {
	dbConn := db.GetConnection(cfg)

	handler := handlers.NewDataHandler(dbConn)
	g := InitGinEngine(cfg, logger)
	group := g.Group(apiGroup(cfg))
	AddHealthCheckRoute(group)
	group.POST("/registerProvider", handler.RegisterAWS)
	return g
}

func (s Server) run() {
	go func() {
		s.logger.Info().Msgf("listning at %s...", s.srv.Addr)

		if err := s.srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				s.logger.Fatal().Err(err).Msg("error starting server")
			}
		}
	}()
}

func apiGroup(cfg config.ServiceConfig) string {
	if !cfg.IsLocalProfile() {
		return fmt.Sprintf("/%s/", cfg.ServiceID)
	}
	return fmt.Sprintf("/")
}

func (s Server) shutDown() {
	s.logger.Info().Msg("Shutting down server...")

	timeout := time.Duration(s.shutdownTimeOut)

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Fatal().Err(err).Msg("Exit server")
	}

	s.logger.Info().Msg("Server exited")
}

func waitForInterrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (s Server) ListenTillDeath(cfg config.ServiceConfig) {
	if !cfg.IsLocalProfile() {
		defer tracer.Stop()
	}
	s.run()
	waitForInterrupt()
	s.shutDown()
}

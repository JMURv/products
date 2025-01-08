package main

import (
	"context"
	"fmt"
	"github.com/JMURv/par-pro/products/internal/cache/redis"
	ctrl "github.com/JMURv/par-pro/products/internal/ctrl"
	sso_ctrl_grpc "github.com/JMURv/par-pro/products/internal/ctrl/sso"
	discovery "github.com/JMURv/par-pro/products/internal/discovery/JMURv/grpc"
	//handler "github.com/JMURv/par-pro/products/internal/handler/http"
	handler "github.com/JMURv/par-pro/products/internal/hdl/http"
	tracing "github.com/JMURv/par-pro/products/internal/metrics/jaeger"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	db "github.com/JMURv/par-pro/products/internal/repo/db"
	cfg "github.com/JMURv/par-pro/products/pkg/config"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "local.config.yaml"

func mustRegisterLogger(mode string) {
	switch mode {
	case "prod":
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	case "dev":
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Panic("panic occurred", zap.Any("error", err))
			os.Exit(1)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	conf := cfg.MustLoad(configPath)
	mustRegisterLogger(conf.Server.Mode)

	go metrics.New(conf.Server.Port + 5).Start(ctx)
	go tracing.Start(ctx, conf.ServiceName, conf.Jaeger)

	dsc := discovery.New(conf.SrvDiscovery, conf.ServiceName, conf.Server)
	if err := dsc.Register(ctx); err != nil {
		zap.L().Fatal("Error registering service", zap.Error(err))
	}

	cache := redis.New(conf.Redis)
	repo := db.New(conf.DB)

	ssoCtrl := sso_ctrl_grpc.New(dsc)
	svc := ctrl.New(repo, cache)
	h := handler.New(svc, ssoCtrl)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c

		zap.L().Info("Shutting down gracefully...")

		cache.Close()
		if err := dsc.Deregister(ctx); err != nil {
			zap.L().Debug("Error deregistering service", zap.Error(err))
		}
		if err := h.Close(); err != nil {
			zap.L().Debug("Error closing handler", zap.Error(err))
		}

		cancel()
		os.Exit(0)
	}()

	zap.L().Info(
		fmt.Sprintf("Starting server on %v://%v:%v", conf.Server.Scheme, conf.Server.Domain, conf.Server.Port),
	)
	h.Start(conf.Server.Port)
}

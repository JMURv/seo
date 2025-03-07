package main

import (
	"context"
	"github.com/JMURv/seo/internal/cache/redis"
	"github.com/JMURv/seo/internal/config"
	"github.com/JMURv/seo/internal/ctrl"
	"github.com/JMURv/seo/internal/ctrl/sso"
	"github.com/JMURv/seo/internal/hdl/http"
	"github.com/JMURv/seo/internal/observability/metrics/prometheus"
	"github.com/JMURv/seo/internal/observability/tracing/jaeger"
	"github.com/JMURv/seo/internal/repo/db"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "configs/local.config.yaml"

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
	defer cancel()

	conf := config.MustLoad(configPath)
	mustRegisterLogger(conf.Mode)

	go prometheus.New(conf.Server.Port + 5).Start(ctx)
	go jaeger.Start(ctx, conf.ServiceName, conf.Jaeger)

	cache := redis.New(conf.Redis)
	repo := db.New(conf.DB)
	svc := ctrl.New(repo, cache)
	h := http.New(svc, sso.New(conf.Services))

	go h.Start(conf.Server.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	zap.L().Info("Shutting down gracefully...")
	if err := h.Close(ctx); err != nil {
		zap.L().Warn("Error closing handler", zap.Error(err))
	}

	if err := cache.Close(); err != nil {
		zap.L().Warn("Failed to close connection to Redis: ", zap.Error(err))
	}

	if err := repo.Close(); err != nil {
		zap.L().Warn("Error closing repository", zap.Error(err))
	}

	os.Exit(0)
}

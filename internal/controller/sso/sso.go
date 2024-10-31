package sso

import (
	"context"
	pb "github.com/JMURv/protos/par-pro"
	ctrl "github.com/JMURv/seo-svc/internal/controller"
	"github.com/JMURv/seo-svc/internal/discovery"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOSvc interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
	GetIDByToken(ctx context.Context, token string) (string, error)
}

type SSO struct {
	discovery *discovery.Discovery
}

func New(discovery *discovery.Discovery) *SSO {
	return &SSO{
		discovery: discovery,
	}
}

func (s *SSO) ValidateToken(ctx context.Context, token string) (bool, error) {
	const op = "products.ValidateToken.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "sso")
	if err != nil {
		zap.L().Debug("failed to find svc", zap.Error(err), zap.String("op", op))
		return false, ctrl.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return false, ctrl.ErrCreateClient
	}
	defer cli.Close()

	res, err := pb.NewSSOClient(cli).ValidateToken(
		ctx, &pb.StringSSOMsg{
			String_: token,
		},
	)
	if err != nil {
		return false, err
	}

	return res.Bool, nil
}

func (s *SSO) GetIDByToken(ctx context.Context, token string) (string, error) {
	const op = "products.ValidateToken.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "sso")
	if err != nil {
		zap.L().Debug("failed to find svc", zap.Error(err), zap.String("op", op))
		return "", ctrl.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return "", ctrl.ErrCreateClient
	}
	defer cli.Close()

	res, err := pb.NewSSOClient(cli).GetUserByToken(
		ctx, &pb.StringSSOMsg{
			String_: token,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

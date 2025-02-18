package sso

import (
	"context"
	"fmt"
	pb "github.com/JMURv/protos/par-pro"
	"github.com/JMURv/seo/internal/config"
	ctrl "github.com/JMURv/seo/internal/ctrl"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOSvc interface {
	ParseClaims(ctx context.Context, token string) (string, error)
}

type SSO struct {
	url string
}

func New(conf *config.ServicesConfig) *SSO {
	return &SSO{
		url: fmt.Sprintf("%v:%v", conf.SSO.Domain, conf.SSO.Port),
	}
}

func (s *SSO) ParseClaims(ctx context.Context, token string) (string, error) {
	const op = "sso.ParseClaims.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	cli, err := grpc.NewClient(s.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return "", ctrl.ErrCreateClient
	}
	defer func(cli *grpc.ClientConn) {
		if err := cli.Close(); err != nil {
			zap.L().Debug(
				"failed to close client",
				zap.String("op", op),
				zap.Error(err),
			)
		}
	}(cli)

	res, err := pb.NewSSOClient(cli).ParseClaims(
		ctx, &pb.SSO_StringMsg{
			String_: token,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Token, nil
}

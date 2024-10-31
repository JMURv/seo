package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/JMURv/seo-svc/api/pb"
	"github.com/JMURv/seo-svc/internal/controller/sso"
	"github.com/JMURv/seo-svc/internal/handler"
	metrics "github.com/JMURv/seo-svc/internal/metrics/prometheus"
	pm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Handler struct {
	pb.SEOServer
	pb.PageServer
	srv  *grpc.Server
	hsrv *health.Server
	ctrl handler.SEOCtrl
	sso  sso.SSOSvc
}

func New(ctrl handler.SEOCtrl, sso sso.SSOSvc) *Handler {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			AuthUnaryInterceptor(sso),
			metrics.SrvMetrics.UnaryServerInterceptor(pm.WithExemplarFromContext(metrics.Exemplar)),
		),
		grpc.ChainStreamInterceptor(
			metrics.SrvMetrics.StreamServerInterceptor(pm.WithExemplarFromContext(metrics.Exemplar)),
		),
	)
	hsrv := health.NewServer()
	hsrv.SetServingStatus("seo", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(srv)
	return &Handler{
		ctrl: ctrl,
		srv:  srv,
		hsrv: hsrv,
		sso:  sso,
	}
}

func (h *Handler) Start(port int) {
	pb.RegisterSEOServer(h.srv, h)
	pb.RegisterPageServer(h.srv, h)
	grpc_health_v1.RegisterHealthServer(h.srv, h.hsrv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := h.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Fatal(err)
	}
}

func (h *Handler) Close() error {
	h.srv.GracefulStop()
	return nil
}

func AuthUnaryInterceptor(sso sso.SSOSvc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			zap.L().Debug("missing metadata")
			return handler(ctx, req)
		}

		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			zap.L().Debug("missing authorization token")
			return handler(ctx, req)
		}

		tokenStr := authHeaders[0]
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		uid, err := sso.GetIDByToken(ctx, tokenStr)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "uid", uid)
		return handler(ctx, req)
	}
}

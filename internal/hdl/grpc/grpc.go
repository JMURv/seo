package grpc

import (
	"errors"
	"fmt"
	"github.com/JMURv/seo/api/grpc/v1/gen"
	"github.com/JMURv/seo/internal/ctrl"
	"github.com/JMURv/seo/internal/ctrl/sso"
	"github.com/JMURv/seo/internal/hdl/grpc/interceptors"
	metrics "github.com/JMURv/seo/internal/observability/metrics/prometheus"
	pm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

type Handler struct {
	gen.SEOServer
	gen.PageServer
	srv  *grpc.Server
	hsrv *health.Server
	ctrl ctrl.AppCtrl
}

func New(name string, ctrl ctrl.AppCtrl, sso sso.SSOSvc) *Handler {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthUnaryInterceptor(sso),
			metrics.SrvMetrics.UnaryServerInterceptor(
				pm.WithExemplarFromContext(metrics.Exemplar),
			),
		),
		grpc.ChainStreamInterceptor(
			metrics.SrvMetrics.StreamServerInterceptor(
				pm.WithExemplarFromContext(metrics.Exemplar),
			),
		),
	)

	reflection.Register(srv)
	hsrv := health.NewServer()
	hsrv.SetServingStatus(name, grpc_health_v1.HealthCheckResponse_SERVING)
	return &Handler{
		ctrl: ctrl,
		srv:  srv,
		hsrv: hsrv,
	}
}

func (h *Handler) Start(port int) {
	gen.RegisterSEOServer(h.srv, h)
	gen.RegisterPageServer(h.srv, h)
	grpc_health_v1.RegisterHealthServer(h.srv, h.hsrv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		zap.L().Fatal("failed to listen", zap.Error(err))
	}

	zap.L().Info(
		"Starting GRPC server",
		zap.String("addr", lis.Addr().String()),
	)
	if err = h.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		zap.L().Fatal("failed to serve", zap.Error(err))
	}
}

func (h *Handler) Close() error {
	h.srv.GracefulStop()
	return nil
}

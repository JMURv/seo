package grpc

import (
	"fmt"
	pb "github.com/JMURv/seo-svc/api/pb"
	"github.com/JMURv/seo-svc/internal/handler"
	metrics "github.com/JMURv/seo-svc/internal/metrics/prometheus"
	pm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
}

func New(ctrl handler.SEOCtrl) *Handler {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
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

	log.Fatal(h.srv.Serve(lis))
}

func (h *Handler) Close() error {
	h.srv.GracefulStop()
	return nil
}

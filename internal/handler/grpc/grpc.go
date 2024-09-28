package grpc

import (
	"fmt"
	pb "github.com/JMURv/par-pro-seo/api/pb"
	ctrl "github.com/JMURv/par-pro-seo/internal/controller"
	metrics "github.com/JMURv/par-pro-seo/internal/metrics/prometheus"
	pm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Handler struct {
	pb.SEOServer
	srv  *grpc.Server
	ctrl *ctrl.Controller
}

func New(ctrl *ctrl.Controller) *Handler {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metrics.SrvMetrics.UnaryServerInterceptor(pm.WithExemplarFromContext(metrics.Exemplar)),
		),
		grpc.ChainStreamInterceptor(
			metrics.SrvMetrics.StreamServerInterceptor(pm.WithExemplarFromContext(metrics.Exemplar)),
		),
	)

	reflection.Register(srv)
	return &Handler{
		ctrl: ctrl,
		srv:  srv,
	}
}

func (h *Handler) Start(port int) {
	pb.RegisterSEOServer(h.srv, h)

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

package grpc

import (
	"fmt"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl/sso"
	"github.com/JMURv/par-pro/products/internal/hdl"
	"github.com/JMURv/par-pro/products/internal/hdl/grpc/interceptors"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	pm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Handler struct {
	pb.ItemServer
	pb.CategoryServer
	pb.PromotionServer
	pb.FavoriteServer
	srv  *grpc.Server
	hsrv *health.Server
	ctrl hdl.Ctrl
	sso  sso.SSOSvc
}

func New(ctrl hdl.Ctrl, sso sso.SSOSvc) *Handler {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthUnaryInterceptor(sso),
			metrics.SrvMetrics.UnaryServerInterceptor(pm.WithExemplarFromContext(metrics.Exemplar)),
		),
		grpc.ChainStreamInterceptor(
			metrics.SrvMetrics.StreamServerInterceptor(pm.WithExemplarFromContext(metrics.Exemplar)),
		),
	)
	hsrv := health.NewServer()
	hsrv.SetServingStatus("products", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(srv)
	return &Handler{
		srv:  srv,
		hsrv: hsrv,
		ctrl: ctrl,
		sso:  sso,
	}
}

func (h *Handler) Start(port int) {
	pb.RegisterItemServer(h.srv, h)
	pb.RegisterCategoryServer(h.srv, h)
	pb.RegisterPromotionServer(h.srv, h)
	pb.RegisterFavoriteServer(h.srv, h)
	grpc_health_v1.RegisterHealthServer(h.srv, h.hsrv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Fatal(h.srv.Serve(lis))
}

func (h *Handler) Close() error {
	h.srv.GracefulStop()
	h.hsrv.Shutdown()
	return nil
}

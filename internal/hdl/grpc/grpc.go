package grpc

import (
	"fmt"
	pb "github.com/JMURv/par-pro/products/api/pb"
	"github.com/JMURv/par-pro/products/internal/ctrl"
	metrics "github.com/JMURv/par-pro/products/internal/metrics/prometheus"
	pm "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

//type Ctrl interface {
//}

type Handler struct {
	pb.ItemServer
	pb.CategoryServer
	pb.PromotionServer
	pb.FavoriteServer
	srv  *grpc.Server
	ctrl ctrl.Controller
}

func New(ctrl ctrl.Controller) *Handler {
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
	pb.RegisterItemServer(h.srv, h)
	pb.RegisterCategoryServer(h.srv, h)
	pb.RegisterPromotionServer(h.srv, h)
	pb.RegisterFavoriteServer(h.srv, h)

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

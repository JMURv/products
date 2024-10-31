package grpc

import (
	"context"
	errs "github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/internal/discovery"
	pb "github.com/JMURv/protos/par-pro"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SSOCtrl struct {
	discovery *discovery.Discovery
}

func New(discovery *discovery.Discovery) *SSOCtrl {
	return &SSOCtrl{
		discovery: discovery,
	}
}

func (s *SSOCtrl) ValidateToken(ctx context.Context, token string) (bool, error) {
	const op = "products.ValidateToken.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "sso")
	if err != nil {
		zap.L().Debug("failed to find svc", zap.Error(err), zap.String("op", op))
		return false, errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return false, errs.ErrCreateClient
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

func (s *SSOCtrl) GetIDByToken(ctx context.Context, token string) (string, error) {
	const op = "products.ValidateToken.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "sso")
	if err != nil {
		zap.L().Debug("failed to find svc", zap.Error(err), zap.String("op", op))
		return "", errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return "", errs.ErrCreateClient
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
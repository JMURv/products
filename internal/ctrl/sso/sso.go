package sso

import (
	"context"
	"errors"
	discovery "github.com/JMURv/par-pro/products/internal/discovery"
	pb "github.com/JMURv/protos/par-pro"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrNotFoundSvc = errors.New("service not found")
var ErrCreateClient = errors.New("failed to create client")

type SSOSvc interface {
	ParseClaims(ctx context.Context, token string) (string, error)
	CreateUser(ctx context.Context, name, email, password string) (string, error)
}

type SSO struct {
	discovery discovery.ServiceDiscovery
}

func New(discovery discovery.ServiceDiscovery) *SSO {
	return &SSO{
		discovery: discovery,
	}
}

func (s *SSO) CreateUser(ctx context.Context, name, email, password string) (string, error) {
	const op = "sso.CreateUser.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "sso")
	if err != nil {
		zap.L().Debug("failed to find svc", zap.Error(err), zap.String("op", op))
		return "", ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return "", ErrCreateClient
	}
	defer cli.Close()

	res, err := pb.NewUsersClient(cli).CreateUser(
		ctx, &pb.SSO_CreateUserReq{
			Name:     name,
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		return "", err
	}

	return res.Uid, nil
}

func (s *SSO) ParseClaims(ctx context.Context, token string) (string, error) {
	const op = "sso.ParseClaims.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "sso")
	if err != nil {
		zap.L().Debug("failed to find svc", zap.Error(err), zap.String("op", op))
		return "", ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return "", ErrCreateClient
	}
	defer cli.Close()

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

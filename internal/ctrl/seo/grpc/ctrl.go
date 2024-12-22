package seo_ctrl_grpc

import (
	"context"
	errs "github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/internal/discovery"
	"github.com/JMURv/par-pro/products/pkg/model/seo"
	pb "github.com/JMURv/protos/par-pro"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SEOCtrl struct {
	discovery *discovery.Discovery
}

func New(discovery *discovery.Discovery) *SEOCtrl {
	return &SEOCtrl{
		discovery: discovery,
	}
}

func (s *SEOCtrl) CreateSEO(ctx context.Context, name, pk string, seo *seo.SEO) error {
	const op = "products.CreateSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "seo")
	if err != nil {
		zap.L().Debug("failed to find banner service", zap.Error(err), zap.String("op", op))
		return errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return errs.ErrCreateClient
	}
	defer cli.Close()

	_, err = pb.NewSEOClient(cli).CreateSEO(
		ctx, &pb.SEOMsg{
			Title:         seo.Title,
			Description:   seo.Description,
			Keywords:      seo.Keywords,
			OGTitle:       seo.OGTitle,
			OGDescription: seo.OGDescription,
			OGImage:       seo.OGImage,
			ObjName:       name,
			ObjPk:         pk,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SEOCtrl) UpdateSEO(ctx context.Context, name, pk string, seo *seo.SEO) error {
	const op = "products.UpdateSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "seo")
	if err != nil {
		zap.L().Debug("failed to find banner service", zap.Error(err), zap.String("op", op))
		return errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return errs.ErrCreateClient
	}
	defer cli.Close()

	_, err = pb.NewSEOClient(cli).UpdateSEO(
		ctx, &pb.SEOMsg{
			Title:         seo.Title,
			Description:   seo.Description,
			Keywords:      seo.Keywords,
			OGTitle:       seo.OGTitle,
			OGDescription: seo.OGDescription,
			OGImage:       seo.OGImage,
			ObjName:       name,
			ObjPk:         pk,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *SEOCtrl) DeleteSEO(ctx context.Context, name, pk string) error {
	const op = "products.DeleteSEO.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "seo")
	if err != nil {
		zap.L().Debug("failed to find banner service", zap.Error(err), zap.String("op", op))
		return errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return errs.ErrCreateClient
	}
	defer cli.Close()

	_, err = pb.NewSEOClient(cli).DeleteSEO(
		ctx, &pb.GetSEOReq{
			Name: name,
			Pk:   pk,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

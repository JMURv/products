package etc_ctrl_grpc

import (
	"context"
	errs "github.com/JMURv/par-pro/products/internal/ctrl"
	"github.com/JMURv/par-pro/products/internal/discovery"
	"github.com/JMURv/par-pro/products/pkg/model/etc"
	pb "github.com/JMURv/protos/par-pro"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EtcCtrl struct {
	discovery *discovery.Discovery
}

func New(discovery *discovery.Discovery) *EtcCtrl {
	return &EtcCtrl{
		discovery: discovery,
	}
}

func slidesToProto(req []etc.BannerSlide) []*pb.SlideMsg {
	res := make([]*pb.SlideMsg, 0, len(req))
	for _, slide := range req {
		res = append(res, &pb.SlideMsg{
			Id:          slide.ID,
			Title:       slide.Title,
			Description: slide.Description,
			Src:         slide.Src,
			Alt:         slide.Alt,
			ButtonText:  slide.ButtonText,
			ButtonHref:  slide.ButtonHref,
			BannerId:    slide.BannerID,
			CreatedAt:   timestamppb.New(slide.CreatedAt),
			UpdatedAt:   timestamppb.New(slide.UpdatedAt),
		})
	}

	return res
}

func (s *EtcCtrl) CreateBanner(ctx context.Context, name, pk string, req *etc.Banner) error {
	const op = "products.CreateBanner.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "etc")
	if err != nil {
		zap.L().Debug("failed to find service", zap.Error(err), zap.String("op", op))
		return errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return errs.ErrCreateClient
	}
	defer cli.Close()

	_, err = pb.NewBannerClient(cli).CreateBanner(ctx, &pb.BannerMsg{
		ObjName: name,
		ObjPk:   pk,
		Slides:  slidesToProto(req.Slides),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *EtcCtrl) UpdateBanner(ctx context.Context, name, pk string, req *etc.Banner) error {
	const op = "products.UpdateBanner.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "etc")
	if err != nil {
		zap.L().Debug("failed to find service", zap.Error(err), zap.String("op", op))
		return errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return errs.ErrCreateClient
	}
	defer cli.Close()

	_, err = pb.NewBannerClient(cli).UpdateBanner(ctx, &pb.CreateAndUpdateBannerReq{
		Name: name,
		Pk:   pk,
		Banner: &pb.BannerMsg{
			Id:      req.ID,
			ObjName: req.OBJName,
			ObjPk:   req.OBJPK,
			Slides:  slidesToProto(req.Slides),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *EtcCtrl) DeleteBanner(ctx context.Context, name, pk string) error {
	const op = "products.DeleteBanner.ctrl"
	span, _ := opentracing.StartSpanFromContext(ctx, op)
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	url, err := s.discovery.FindServiceByName(ctx, "etc")
	if err != nil {
		zap.L().Debug("failed to find service", zap.Error(err), zap.String("op", op))
		return errs.ErrNotFoundSvc
	}

	cli, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.L().Debug("failed to create client", zap.Error(err), zap.String("op", op))
		return errs.ErrCreateClient
	}
	defer cli.Close()

	_, err = pb.NewBannerClient(cli).DeleteBanner(ctx, &pb.GetBannerReq{
		Name: name,
		Pk:   pk,
	})
	if err != nil {
		return err
	}
	return nil
}

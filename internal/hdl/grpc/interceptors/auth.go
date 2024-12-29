package interceptors

import (
	"context"
	"github.com/JMURv/par-pro/products/internal/ctrl/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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

		uid, err := sso.ParseClaims(ctx, tokenStr)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, "uid", uid)
		return handler(ctx, req)
	}
}

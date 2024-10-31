package sso

import (
	"context"
)

type SSOCtrl interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
	GetIDByToken(ctx context.Context, token string) (string, error)
}

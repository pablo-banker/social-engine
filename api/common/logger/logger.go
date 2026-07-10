package logger

import (
	"context"

	"go.uber.org/zap"
)

func L(ctx context.Context) *zap.Logger {
	clientId := clientExecIDFromCtx(ctx)
	userId := userIDFromCtx(ctx)

	return zap.L().With(
		zap.String("user_id", userId),
		zap.String("support_code", clientId),
	)
}

func clientExecIDFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value("support_code"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func userIDFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value("user_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

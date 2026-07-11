package logger

import (
	"context"

	"go.uber.org/zap"
)

func L(ctx context.Context) *zap.Logger {
	return zap.L().With(
		zap.String("trace_id", traceIDFromCtx(ctx)),
		zap.String("user_id", userIDFromCtx(ctx)),
		zap.String("support_code", clientExecIDFromCtx(ctx)),
	)
}

func traceIDFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v := ctx.Value("trace_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
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

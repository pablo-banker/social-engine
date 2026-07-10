package logger

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func Test_ContextHelpersAndLogger(t *testing.T) {
	// garante que zap.L() exista e não dê nil
	core, obs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)

	t.Run("clientExecIDFromCtx", func(t *testing.T) {
		t.Run("nil ctx -> empty", func(t *testing.T) {
			assert.Equal(t, "", clientExecIDFromCtx(context.TODO()))
		})

		t.Run("missing key -> empty", func(t *testing.T) {
			ctx := context.Background()
			assert.Equal(t, "", clientExecIDFromCtx(ctx))
		})

		t.Run("wrong type -> empty", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "support_code", 123)
			assert.Equal(t, "", clientExecIDFromCtx(ctx))
		})

		t.Run("success", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "support_code", "SC-123")
			assert.Equal(t, "SC-123", clientExecIDFromCtx(ctx))
		})
	})

	t.Run("userIDFromCtx", func(t *testing.T) {
		t.Run("nil ctx -> empty", func(t *testing.T) {
			assert.Equal(t, "", userIDFromCtx(context.TODO()))
		})

		t.Run("missing key -> empty", func(t *testing.T) {
			ctx := context.Background()
			assert.Equal(t, "", userIDFromCtx(ctx))
		})

		t.Run("wrong type -> empty", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "user_id", uuid.New())
			assert.Equal(t, "", userIDFromCtx(ctx))
		})

		t.Run("success", func(t *testing.T) {
			ctx := context.WithValue(context.Background(), "user_id", "u_123")
			assert.Equal(t, "u_123", userIDFromCtx(ctx))
		})
	})

	t.Run("L", func(t *testing.T) {
		t.Run("ctx nil -> still returns logger with fields (empty values ok)", func(t *testing.T) {

			l := L(context.TODO())
			assert.NotNil(t, l)

			l.Info("hello")
			logs := obs.TakeAll()
			assert.Len(t, logs, 1)

			fields := map[string]any{}
			for _, f := range logs[0].Context {
				fields[f.Key] = f.Interface
				if f.Interface == nil {
					// para strings, o observer usa String e deixa em .String
					fields[f.Key] = f.String
				}
			}

			// os campos existem
			assert.Contains(t, fields, "trace_id")
			assert.Contains(t, fields, "user_id")
			assert.Contains(t, fields, "support_code")
		})

		t.Run("ctx with values -> attaches user_id and support_code", func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "support_code", "SUP-999")
			ctx = context.WithValue(ctx, "user_id", "user-abc")

			l := L(ctx)
			assert.NotNil(t, l)

			l.Info("x")
			logs := obs.TakeAll()
			assert.Len(t, logs, 1)

			got := map[string]string{}
			for _, f := range logs[0].Context {
				// no observer, string fica em f.String
				got[f.Key] = f.String
			}

			assert.Equal(t, "user-abc", got["user_id"])
			assert.Equal(t, "SUP-999", got["support_code"])
			assert.Contains(t, got, "trace_id")
		})

		t.Run("ctx wrong types -> empty strings for user_id/support_code", func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "support_code", 999) // wrong type
			ctx = context.WithValue(ctx, "user_id", true)     // wrong type

			l := L(ctx)
			assert.NotNil(t, l)

			l.Info("y")
			logs := obs.TakeAll()
			assert.Len(t, logs, 1)

			got := map[string]string{}
			for _, f := range logs[0].Context {
				got[f.Key] = f.String
			}

			assert.Equal(t, "", got["user_id"])
			assert.Equal(t, "", got["support_code"])
			assert.Contains(t, got, "trace_id")
		})
	})
}

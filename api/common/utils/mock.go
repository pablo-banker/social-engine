package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"social-engine/common/logger"
	"social-engine/common/repositories"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type TestFiberReqConfig struct {
	Body             any
	MockData         any
	Headers          map[string]string
	QueryParams      map[string]any
	QueryParamsOrder []string
	URLParams        map[string]string
	Middlewares      []fiber.Handler
	Handler          fiber.Handler
	CtxValues        map[any]any
}

func (t *TestFiberReqConfig) GetQueryString() string {
	if len(t.QueryParams) == 0 {
		return ""
	}

	if len(t.QueryParamsOrder) > 0 {
		params := make([]string, 0, len(t.QueryParams))
		for _, key := range t.QueryParamsOrder {
			if val, ok := t.QueryParams[key]; ok {
				params = append(params, key+"="+url.QueryEscape(fmt.Sprintf("%v", val)))
			}
		}
		return "?" + strings.Join(params, "&")
	}

	values := url.Values{}
	for key, val := range t.QueryParams {
		values.Set(key, fmt.Sprintf("%v", val))
	}
	return "?" + values.Encode()
}

func NewTestCtx(config TestFiberReqConfig) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, repositories.GormTestContext, config.MockData)

	if config.CtxValues != nil {
		for key, value := range config.CtxValues {
			ctx = context.WithValue(ctx, key, value)
		}
	}

	return ctx
}

func NewTestFiberCtx(config TestFiberReqConfig) *fiber.Ctx {
	app := fiber.New()
	ctx := NewTestCtx(config)

	// Create a dummy fasthttp.RequestCtx for testing
	reqCtx := &fasthttp.RequestCtx{}

	// Acquire a fiber.Ctx instance from the app
	c := app.AcquireCtx(reqCtx)
	c.SetUserContext(ctx)
	c.Locals("userId", uuid.Nil)

	return c
}

func NewTestFiberReq(config TestFiberReqConfig) (int, map[string]any, error) {
	app := fiber.New()

	// Mock Context
	ctx := NewTestCtx(config)

	// Monta Body
	var bodyBytes []byte
	if config.Body != nil {
		var err error
		bodyBytes, err = json.Marshal(config.Body)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	// Monta URL Path e Route Pattern
	urlPath := "/test"
	routePattern := "/test"

	if len(config.URLParams) > 0 {
		for key := range config.URLParams {
			routePattern += "/:" + key
		}
		for _, value := range config.URLParams {
			urlPath += "/" + value
		}
	}

	// Monta Query String
	queryString := config.GetQueryString()

	// Registra Middlewares
	for _, m := range config.Middlewares {
		app.Use(m)
	}

	// Registra Handler
	if config.Handler != nil {
		app.Post(routePattern, func(c *fiber.Ctx) error {
			// Injeta o userContext no contexto da requisição
			c.SetUserContext(ctx)
			return config.Handler(c)
		})
	}

	//  Cria request real
	req := httptest.NewRequest(fiber.MethodPost, urlPath+queryString, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	//  Adiciona Headers
	if config.Headers != nil {
		for key, value := range config.Headers {
			req.Header.Set(key, value)
		}
	}

	//  Dispara o ciclo real
	res, err := app.Test(req, -1)
	if err != nil {
		return 0, nil, fmt.Errorf("fiber app test error: %w", err)
	}

	//  Faz o parse da resposta (JSON)
	var parsedBody map[string]interface{}
	if res.Body != nil {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return res.StatusCode, nil, fmt.Errorf("failed to read response body: %w", err)
		}

		cleanData := strings.ReplaceAll(string(data), "\\", "")
		trimmedData := strings.TrimSpace(cleanData)

		if trimmedData == "" || trimmedData == `""` {
			return res.StatusCode, nil, errors.New("request failed")
		}

		if err = json.Unmarshal(data, &parsedBody); err != nil {
			logger.L(ctx).Error("Failed to unmarshal error response", zap.Error(err))
			return res.StatusCode, nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return res.StatusCode, parsedBody, nil
}

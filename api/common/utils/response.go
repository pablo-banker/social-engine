package utils

import (
	"fmt"
	"social-engine/common/apiErrors"
	"social-engine/common/logger"
	"social-engine/common/repositories"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Response[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Status      int    `json:"status"`
	Code        string `json:"code"`
	Error       string `json:"error"`
	Details     string `json:"details,omitempty"`
	SupportCode string `json:"supportCode,omitempty"`
}

type SuccessResponse struct {
	Message string                 `json:"message,omitempty"`
	Status  string                 `json:"status,omitempty"`
	ID      string                 `json:"id,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func BuildErrorResponse(c *fiber.Ctx, err *apiErrors.APIError) error {
	resp := ErrorResponse{
		Status:  err.HTTPStatus,
		Code:    err.Code,
		Error:   err.Error(),
		Details: err.Details,
	}

	if err.Internal != "" {
		logger.L(c.UserContext()).Error(err.Message, zap.String("code", err.Code), zap.String("internal", err.Internal))
	}

	testControl := c.UserContext().Value(repositories.GormTestContext)
	if testControl != nil {
		fmt.Printf("Error: %s, Code: %s, Internal: %s, Details: %s\n", err.Message, err.Code, err.Internal, err.Details)
	}

	data := c.UserContext().Value("support_code")
	if data != nil {
		if supportCode, ok := data.(string); ok {
			resp.SupportCode = supportCode
		}
	}

	return c.Status(err.HTTPStatus).JSON(resp)
}

func BuildSuccessResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(Response[interface{}]{
		Data: data,
	})
}

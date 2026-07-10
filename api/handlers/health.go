package handlers

import (
	"net/http"
	"social-engine/common/repositories"
	"social-engine/common/utils"

	"github.com/gofiber/fiber/v2"
)

type HealthResponse struct {
	Status string `json:"status"`

	Checks []HealthCheck `json:"checks,omitempty"`
}

type HealthCheck struct {
	Name  string `json:"name"`
	Alive bool   `json:"alive"`
}

// Health is a handler that checks the health of the API
// @Summary Check the health of the API
// @Description Check the health of the API
// @Tags Default
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func Health(c *fiber.Ctx) error {
	var (
		ctx  = c.UserContext()
		repo = repositories.InitializeRepoInstance(ctx)
	)

	dbUp := HealthCheck{
		Name:  "Database",
		Alive: repo.Ping(ctx) == nil,
	}

	status := "UP"
	code := http.StatusOK

	if !dbUp.Alive {
		status = "DOWN"
		code = http.StatusServiceUnavailable
	}

	return utils.BuildSuccessResponse(c, code, HealthResponse{
		Status: status,
		Checks: []HealthCheck{dbUp},
	})
}

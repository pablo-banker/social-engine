package handlers

import (
	"fmt"
	"net/http"
	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HealthCheck(t *testing.T) {
	cases := []struct {
		name              string
		mockData          []repositories.MockPayload
		wantStatus        int
		wantServiceStatus string
	}{
		{
			name: "Health check - all services up",
			mockData: []repositories.MockPayload{
				{
					Type:          constants.RepositoryPing,
					ExpectedError: nil,
				},
			},
			wantStatus:        http.StatusOK,
			wantServiceStatus: "UP",
		},
		{
			name: "Health check - database down",
			mockData: []repositories.MockPayload{
				{
					Type:          constants.RepositoryPing,
					ExpectedError: fmt.Errorf("database down"),
				},
			},
			wantStatus:        http.StatusServiceUnavailable,
			wantServiceStatus: "DOWN",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
				Handler:  Health,
				MockData: tt.mockData,
			})
			assert.NoError(t, err)

			data, ok := body["data"].(map[string]any)
			assert.True(t, ok, "data should be a map")

			assert.Equal(t, tt.wantStatus, statusCode)
			assert.Equal(t, tt.wantServiceStatus, data["status"])
		})
	}
}

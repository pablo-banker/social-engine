package handlers

import (
	"net/http"
	"os"
	"testing"

	"social-engine/common/apiErrors"
	"social-engine/common/models"
	"social-engine/common/repositories"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"
	"social-engine/common/utils"
	"social-engine/common/validation"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestMain wires up the validator and a JWT secret for the whole handlers
// package test run (both are configured in main() at runtime).
func TestMain(m *testing.M) {
	validation.InitValidator()
	os.Setenv("JWT_SECRET", "test-secret")
	os.Exit(m.Run())
}

func Test_Register(t *testing.T) {
	savedUser := entities.User{
		Base:     entities.Base{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
		Name:     "Ada Lovelace",
		Username: "ada",
		Email:    "ada@example.com",
		AvatarID: models.DefaultAvatarID,
		BannerID: models.DefaultBannerID,
	}

	cases := []struct {
		name       string
		body       registerRequest
		mockData   []repositories.MockPayload
		wantStatus int
		wantCode   string // apiError code for failures; "" on success
	}{
		{
			name: "success",
			body: registerRequest{FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Password: "secret123"},
			mockData: []repositories.MockPayload{
				{
					Type: constants.RepositoryVerify,
					Params: &entities.QueryParams{
						Query: entities.Query{Filters: "lower(email) = ?", Values: []any{"ada@example.com"}},
					},
					ExpectedResult: false,
				},
				{
					Type: constants.RepositoryVerify,
					Params: &entities.QueryParams{
						Query: entities.Query{Filters: "lower(username) = ?", Values: []any{"ada"}},
					},
					ExpectedResult: false,
				},
				{
					Type:           constants.RepositorySave,
					ExpectedResult: savedUser,
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "email already taken",
			body: registerRequest{FirstName: "Ada", LastName: "Lovelace", Email: "ada@example.com", Password: "secret123"},
			mockData: []repositories.MockPayload{
				{
					Type: constants.RepositoryVerify,
					Params: &entities.QueryParams{
						Query: entities.Query{Filters: "lower(email) = ?", Values: []any{"ada@example.com"}},
					},
					ExpectedResult: true,
				},
			},
			wantStatus: http.StatusConflict,
			wantCode:   apiErrors.ErrEmailTaken.Code,
		},
		{
			name:       "validation error - missing fields",
			body:       registerRequest{Email: "not-an-email"},
			mockData:   []repositories.MockPayload{},
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   apiErrors.ErrValidation.Code,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
				Handler:  Register,
				Body:     tt.body,
				MockData: tt.mockData,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, statusCode)

			if tt.wantCode != "" {
				assert.Equal(t, tt.wantCode, body["code"])
				return
			}

			data, ok := body["data"].(map[string]any)
			assert.True(t, ok, "data should be a map")

			user, ok := data["user"].(map[string]any)
			assert.True(t, ok, "user should be a map")
			assert.Equal(t, "ada@example.com", user["email"])
			assert.Equal(t, "ada", user["username"])
			assert.Equal(t, "Ada Lovelace", user["name"])
			assert.NotEmpty(t, data["accessToken"])
		})
	}
}

func Test_Login(t *testing.T) {
	hash, err := models.HashPassword("secret123")
	assert.NoError(t, err)

	storedUser := entities.User{
		Base:         entities.Base{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222")},
		Name:         "Ada Lovelace",
		Username:     "ada",
		Email:        "ada@example.com",
		AvatarID:     models.DefaultAvatarID,
		PasswordHash: hash,
	}

	emailParams := &entities.QueryParams{
		Query: entities.Query{Filters: "lower(email) = ?", Values: []any{"ada@example.com"}},
	}

	cases := []struct {
		name       string
		body       loginRequest
		mockData   []repositories.MockPayload
		wantStatus int
		wantCode   string
	}{
		{
			name: "success",
			body: loginRequest{Email: "ada@example.com", Password: "secret123"},
			mockData: []repositories.MockPayload{
				{Type: constants.RepositoryFindOne, Params: emailParams, ExpectedResult: storedUser},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			body: loginRequest{Email: "ada@example.com", Password: "wrong-pass"},
			mockData: []repositories.MockPayload{
				{Type: constants.RepositoryFindOne, Params: emailParams, ExpectedResult: storedUser},
			},
			wantStatus: http.StatusUnauthorized,
			wantCode:   apiErrors.ErrInvalidCredentials.Code,
		},
		{
			name: "user not found",
			body: loginRequest{Email: "ada@example.com", Password: "secret123"},
			mockData: []repositories.MockPayload{
				{Type: constants.RepositoryFindOne, Params: emailParams, ExpectedError: gorm.ErrRecordNotFound},
			},
			wantStatus: http.StatusUnauthorized,
			wantCode:   apiErrors.ErrInvalidCredentials.Code,
		},
		{
			name:       "validation error - bad email",
			body:       loginRequest{Email: "nope", Password: "secret123"},
			mockData:   []repositories.MockPayload{},
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   apiErrors.ErrValidation.Code,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body, err := utils.NewTestFiberReq(utils.TestFiberReqConfig{
				Handler:  Login,
				Body:     tt.body,
				MockData: tt.mockData,
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, statusCode)

			if tt.wantCode != "" {
				assert.Equal(t, tt.wantCode, body["code"])
				return
			}

			data, ok := body["data"].(map[string]any)
			assert.True(t, ok, "data should be a map")

			user, ok := data["user"].(map[string]any)
			assert.True(t, ok, "user should be a map")
			assert.Equal(t, "ada@example.com", user["email"])
			assert.NotEmpty(t, data["accessToken"])
		})
	}
}

package http

import (
	"context"
	"encoding/json"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/domain/entities"
	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases/mocks"
	context_keys "eventsguard/internal/infrastructure/config/context"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	testCases := []struct {
		name             string
		expectedStatus   int
		expectedResponse map[string]interface{}
		mockOutput       bool
		mockError        *errors.AppError
	}{
		{
			name:             "Successful logout",
			expectedStatus:   204,
			expectedResponse: nil,
			mockOutput:       true,
			mockError:        nil,
		},
		{
			name:           "Token not found",
			expectedStatus: 404,
			expectedResponse: map[string]interface{}{
				"title":  "Unauthorized",
				"status": 401,
				"detail": "Invalid username/password or user doesn't exist",
			},
			mockOutput: false,
			mockError:  errors.NewNotFoundError("Error deleting token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, api := humatest.New(t)

			// Create a mock implementation of the LogoutUseCase interface.
			logout := auth_usecases_ports.NewMockLogoutUseCase(t)

			// Register the logout handler to the API using the mock.
			NewLoginRouterHandler(api, nil, logout)

			ctx := context.Background()

			// Create a mock token for testing
			token := &entities.Token{
				Token: "mock_token",
			}

			// Set the token in the context
			ctx = context.WithValue(ctx, context_keys.TokenContextKey, token)

			// Set up expectations for the mock use case if applicable.
			if tc.mockOutput || tc.mockError != nil {
				logout.EXPECT().Execute(ctx, token).Return(tc.mockOutput, tc.mockError).Times(1)
			}

			// Test the handler by sending a POST request to the logout endpoint.
			tokenKeader := "Autorization: Bearer " + token.Token
			resp := api.PostCtx(ctx, "/auth/logout", tokenKeader, map[string]any{})

			// Assert that the response status matches the expected status.
			assert.Equal(t, tc.expectedStatus, resp.Result().StatusCode)

			// If an error is expected, validate the error message.
			if tc.expectedStatus != 204 {
				var respBody map[string]interface{}
				if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}

				println(&respBody)

				// // Check that the "errors" field exists.
				// errors, ok := respBody["errors"].([]interface{})
				// assert.True(t, ok, "errors field not found in response")

				// // Validate message.
				// assert.Equal(t, tc.mockError.Message, errors[0].(map[string]interface{})["message"])
			}
		})
	}
}

package http

import (
	"context"
	"encoding/json"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/constants"
	"eventsguard/internal/auth/domain/entities"
	auth_usecases_ports "eventsguard/internal/auth/domain/ports/usecases/mocks"
	"eventsguard/internal/auth/dtos"

	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/danielgtaylor/huma/v2/humatest"
)

const expectedLengthError = "expected length >= 1"

// tokenFixture returns a mock token entity used for testing.
func tokenFixture() *entities.Token {
	return &entities.Token{
		Device:       constants.DeviceApi,
		Token:        "token",
		RefreshToken: "refresh_token",
		User: entities.FakeUser{
			ID:        "id",
			FirstName: "first_name",
			LastName:  nil,
			Username:  "test@test.com",
			IsAdmin:   false,
		},
		UserID:    "test@test.com",
		ExpiresAt: nil,
		CreatedAt: time.Time{},
		Expiracy:  nil,
	}
}

func TestLogin(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedErrors []struct {
			expectedMessage string
			expectedField   string
		}
		mockOutput *entities.Token
		mockError  *errors.AppError
	}{
		{
			name: "Empty username",
			payload: map[string]interface{}{
				"username": "",
				"password": "1234",
			},
			expectedStatus: 422,
			expectedErrors: []struct {
				expectedMessage string
				expectedField   string
			}{
				{
					expectedMessage: expectedLengthError,
					expectedField:   "body.username",
				},
			},
			mockOutput: nil,
			mockError:  nil,
		},
		{
			name: "Empty password",
			payload: map[string]interface{}{
				"username": "daniel",
				"password": "",
			},
			expectedStatus: 422,
			expectedErrors: []struct {
				expectedMessage string
				expectedField   string
			}{
				{
					expectedMessage: expectedLengthError,
					expectedField:   "body.password",
				},
			},
			mockOutput: nil,
			mockError:  nil,
		},
		{
			name: "Empty username and password",
			payload: map[string]interface{}{
				"username": "",
				"password": "",
			},
			expectedStatus: 422,
			expectedErrors: []struct {
				expectedMessage string
				expectedField   string
			}{
				{
					expectedMessage: expectedLengthError,
					expectedField:   "body.username",
				},
				{
					expectedMessage: expectedLengthError,
					expectedField:   "body.password",
				},
			},
			mockOutput: nil,
			mockError:  nil,
		},
		{
			name: "Valid credentials",
			payload: map[string]interface{}{
				"username": "daniel",
				"password": "1234",
			},
			expectedStatus: 200,
			expectedErrors: nil,
			mockOutput:     tokenFixture(),
			mockError:      nil,
		},
		{
			name: "User not found",
			payload: map[string]interface{}{
				"username": "daniel",
				"password": "1234",
			},
			expectedStatus: 401,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      errors.NewAuthError("Invalid username/password or user doesn't exist"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, api := humatest.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock implementation of the LoginUseCase interface.

			login := auth_usecases_ports.NewMockLoginUseCase(t)

			// Register the login handler to the API using the mock.
			NewLoginRouterHandler(api, login, nil)

			ctx := context.Background()

			useCaseInDto := dtos.LoginInputDTO{
				Username: tc.payload["username"].(string),
				Password: tc.payload["password"].(string),
			}

			// Set up expectations for the mock use case if applicable.
			if tc.mockOutput != nil || tc.mockError != nil {
				login.EXPECT().Execute(ctx, useCaseInDto, constants.DeviceApi).Return(tc.mockOutput, tc.mockError).Times(1)
			}

			// Test the handler by sending a POST request to the login endpoint with the payload.
			resp := api.PostCtx(ctx, "/auth/login", tc.payload)

			// Assert that the response status matches the expected status.
			assert.Equal(t, tc.expectedStatus, resp.Result().StatusCode)

			// If an error is expected, validate the error message and location.
			if tc.expectedErrors != nil {
				var respBody map[string]interface{}
				if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}

				// Check that the "errors" field exists.
				errors, ok := respBody["errors"].([]interface{})
				assert.True(t, ok, "errors field not found in response")

				// Validate message and location.
				for _, expectedError := range tc.expectedErrors {
					var found bool
					for _, errItem := range errors {
						if errMap, ok := errItem.(map[string]interface{}); ok {
							if msg, ok := errMap["message"].(string); ok && msg == expectedError.expectedMessage {
								// Check that the "location" matches the expected field.
								if location, ok := errMap["location"].(string); ok && location == expectedError.expectedField {
									found = true
									break
								}
							}
						}
					}
					// Ensure the error was found in the response.
					assert.True(t, found, "Expected error message and field not found")
				}
			}

			// If success is expected, validate the token output.
			if tc.expectedStatus == 200 {
				var token entities.Token
				if err := json.Unmarshal(resp.Body.Bytes(), &token); err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}
				assert.Equal(t, tc.mockOutput, &token)
			}
		})
	}
}

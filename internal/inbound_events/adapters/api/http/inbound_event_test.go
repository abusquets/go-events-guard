package http

import (
	"context"
	"eventsguard/internal/app/errors"
	"eventsguard/internal/auth/domain/entities"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/jaswdr/faker/v2"
)

func TestAddInboundEvent(t *testing.T) {
	eventTypeCorePersonCreated := "core.person.created"

	fake := faker.New()
	// Define test cases
	testCases := []struct {
		name           string
		queryString    string
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
			name: "core.person.created - Success",
			payload: map[string]interface{}{
				"type": eventTypeCorePersonCreated,
				"data": map[string]interface{}{
					"id":        "123",
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": fake.UUID().V4(),
				},
			},
			expectedStatus: 201,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "client_id as path parameter - Success",
			queryString: "client_id_path=data.client_id",
			payload: map[string]interface{}{
				"type": eventTypeCorePersonCreated,
				"data": map[string]interface{}{
					"id":        "123",
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": fake.UUID().V4(),
				},
			},
			expectedStatus: 201,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "client_id as path parameter - Unsuccessful, client_id path not found",
			queryString: "client_id_path=data.client_id",
			payload: map[string]interface{}{
				"type": "core.person.created",
				"data": map[string]interface{}{
					"id": "123",
				},
			},
			expectedStatus: 422,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "client_id as value parameter - Success",
			queryString: "client_id=1234",
			payload: map[string]interface{}{
				"id":   fake.UUID().V4(),
				"type": eventTypeCorePersonCreated,
				"data": map[string]interface{}{
					"id":        "123",
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": "1234",
				},
			},
			expectedStatus: 201,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "type path parameter - Success",
			queryString: "type_path=mytype",
			payload: map[string]interface{}{
				"mytype": eventTypeCorePersonCreated,
				"data": map[string]interface{}{
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": fake.UUID().V4(),
				},
			},
			expectedStatus: 201,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "type path parameter - UnSuccess",
			queryString: "type_path=type",
			payload: map[string]interface{}{
				"mytype": eventTypeCorePersonCreated,
				"data": map[string]interface{}{
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": fake.UUID().V4(),
				},
			},
			expectedStatus: 422,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "payload path parameter - Success",
			queryString: "type_path=payload.type&client_id_path=payload.client_id&payload_path=payload",
			payload: map[string]interface{}{
				"payload": map[string]interface{}{
					"type":      eventTypeCorePersonCreated,
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": fake.UUID().V4(),
				},
			},
			expectedStatus: 201,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
		{
			name:        "payload path parameter - UnSuccess",
			queryString: "type_path=payload.type&client_id_path=payload.client_id&payload_path=payload",
			payload: map[string]interface{}{
				"data": map[string]interface{}{
					"type":      eventTypeCorePersonCreated,
					"name":      fake.Person().Name(),
					"email":     fake.Person().Contact().Email,
					"client_id": fake.UUID().V4(),
				},
			},
			expectedStatus: 422,
			expectedErrors: nil,
			mockOutput:     nil,
			mockError:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, api := humatest.New(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Register the login handler to the API using the mock.
			NewInboundEventRouterHandler(api)

			ctx := context.Background()

			t.Log("Say bye")

			// // Test the handler by sending a POST request to the login endpoint with the payload.
			url := "/inbound-event?" + tc.queryString
			resp := api.PostCtx(ctx, url, tc.payload)

			// Assert that the response status matches the expected status.
			assert.Equal(t, tc.expectedStatus, resp.Result().StatusCode)

			// // If an error is expected, validate the error message and location.
			// if tc.expectedErrors != nil {
			// 	var respBody map[string]interface{}
			// 	if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
			// 		t.Fatalf("Failed to unmarshal response body: %v", err)
			// 	}

			// 	// Check that the "errors" field exists.
			// 	errors, ok := respBody["errors"].([]interface{})
			// 	assert.True(t, ok, "errors field not found in response")

			// 	// Validate message and location.
			// 	for _, expectedError := range tc.expectedErrors {
			// 		var found bool
			// 		for _, errItem := range errors {
			// 			if errMap, ok := errItem.(map[string]interface{}); ok {
			// 				if msg, ok := errMap["message"].(string); ok && msg == expectedError.expectedMessage {
			// 					// Check that the "location" matches the expected field.
			// 					if location, ok := errMap["location"].(string); ok && location == expectedError.expectedField {
			// 						found = true
			// 						break
			// 					}
			// 				}
			// 			}
			// 		}
			// 		// Ensure the error was found in the response.
			// 		assert.True(t, found, "Expected error message and field not found")
			// 	}
			// }

			// // If success is expected, validate the token output.
			// if tc.expectedStatus == 200 {
			// 	var token entities.Token
			// 	if err := json.Unmarshal(resp.Body.Bytes(), &token); err != nil {
			// 		t.Fatalf("Failed to unmarshal response body: %v", err)
			// 	}
			// 	assert.Equal(t, tc.mockOutput, &token)
			// }
		})
	}
}

package http

import (
	"context"
	"fmt"
	"net/http"

	app_errors "eventsguard/internal/app/errors"
	"eventsguard/internal/inbound_events/dtos"
	"eventsguard/internal/infrastructure/server/server_errors"
	"eventsguard/internal/utils/entities"

	"github.com/PaesslerAG/jsonpath"
	"github.com/danielgtaylor/huma/v2"
)

type InboundEventRouterHandler struct{}

func NewInboundEventRouterHandler(api huma.API) *InboundEventRouterHandler {
	ret := InboundEventRouterHandler{}

	huma.Register(api, huma.Operation{
		OperationID:   "add-inbound-event",
		Method:        http.MethodPost,
		Path:          "/inbound-event",
		Description:   "Add an Inbound Event",
		Tags:          []string{"InboundEvent"},
		DefaultStatus: http.StatusCreated,
		Security: []map[string][]string{
			{"TokenAuth": {"api"}},
		},
	}, handleAddInboundEvent)

	return &ret
}

func handleAddInboundEvent(ctx context.Context, input *dtos.CreateInboundEventRequest) (*dtos.CreateInboundEventResponse, error) {
	payload := input.Body
	clientIDValue, err := extractValueFromPath(input.ClientIDPath, input.ClientIDValue, payload, "Client ID Path")
	if err != nil {
		return nil, err
	}

	typePath := input.TypePath
	if typePath == "" {
		typePath = "type"
	}

	typeValue, err := extractValueFromPath(typePath, input.TypeValue, payload, "Type Path")

	if err != nil {
		return nil, err
	}
	if typeValue == "" {
		return nil, createValidationError("Type Value not found")
	}

	var inPayload map[string]interface{}

	if input.PayloadPath != "" {
		exp := "$." + input.PayloadPath
		result, err := jsonpath.Get(exp, payload)
		if err != nil {
			return nil, createValidationError("Payload Path not found")
		}
		inPayload = result.(map[string]interface{})
	} else {
		inPayload = payload
	}
	fmt.Printf("%+v\n", inPayload)

	inDto := createInboundEventInput(typeValue, clientIDValue, inPayload)

	fmt.Printf("%+v\n", inDto)

	resp := &dtos.CreateInboundEventResponse{
		Body: dtos.CreateInboundEventBody{
			ID: entities.ID("123" + clientIDValue),
		},
	}
	return resp, nil
}

func extractValueFromPath(path, defaultValue string, payload map[string]interface{}, errorMessage string) (string, error) {
	if path == "" {
		return defaultValue, nil
	}
	exp := "$." + path
	result, err := jsonpath.Get(exp, payload)
	if err != nil {
		return "", createValidationError(errorMessage)
	}
	return result.(string), nil
}

func createValidationError(message string) error {
	err := app_errors.NewValidationError(message)
	return server_errors.AppErrorToHumaError(*err)
}

func createInboundEventInput(
	typeValue,
	clientIDValue string,
	payload map[string]interface{},
) dtos.CreateInboundEventInput {
	inDto := dtos.CreateInboundEventInput{
		Type:    typeValue,
		Payload: payload,
	}
	if clientIDValue != "" {
		clientID := entities.ID(clientIDValue)
		inDto.ClientID = &clientID
	}
	return inDto
}

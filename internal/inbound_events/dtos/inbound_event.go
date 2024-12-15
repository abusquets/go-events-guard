package dtos

import (
	"eventsguard/internal/utils/entities"
	"time"
)

// {
//     "id": "unique-event-id",  // Identificador únic de l'esdeveniment
//     "type": "entity.action",  // Tipus d'esdeveniment, per exemple, "core.person.created"
//     "source": "service-name", // Servei o component que genera l'esdeveniment
//     "timestamp": "2024-12-12T10:00:00Z", // Marca temporal UTC
//     "version": "1.0",        // Versió de l'esdeveniment
//     "data": {                // Dades específiques de l'esdeveniment
//         "id": "123",
//         "name": "John Doe",
//         "email": "John@gmail.com"
//     },
//     "metadata": {            // Context addicional (opcional)
//         "traceId": "abc-123",
//         "userId": "456"
//     }
// }

// cloudevents

// {
//     "specversion": "1.0",
//     "id": "A234-1234-1234",
//     "source": "/mycontext",
//     "type": "com.example.someevent",
//     "time": "2024-12-12T10:00:00Z",
//     "data": {
//         "appinfoA": "abc"
//     }
// }

type CreateInboundEventInput struct {
	EmployeeID entities.ID
	Type       string
	ClientID   *entities.ID
	Payload    map[string]interface{}
	SendAt     time.Time
}

type CreateInboundEventRequest struct {
	Body            map[string]interface{} `json:"body" validate:"required"`
	Source          string                 `query:"source" doc:"Source Value"`
	PayloadPath     string                 `query:"payload_path" doc:"Path to Payload in the Body, using JsonPath, default is the Body"`
	ClientIDValue   string                 `query:"client_id" doc:"Client ID Value"`
	ClientIDPath    string                 `query:"client_id_path" doc:"Path to Client ID in the Body, using JsonPath"`
	TypeValue       string                 `query:"type" doc:"Type Value"`
	TypePath        string                 `query:"type_path" doc:"Path to Type in the Body, using JsonPath"`
	CreatedAt       string                 `query:"created_at" doc:"Created At Value"`
	CreatedAtPath   string                 `query:"created_at_path" doc:"Path to CreatedAt in the Body, using JsonPath"`
	CreatedAtFormat string                 `query:"created_at_format" doc:"Created At Format"`
}

type CreateInboundEventBody struct {
	ID entities.ID
}

type CreateInboundEventResponse struct {
	Body CreateInboundEventBody
}

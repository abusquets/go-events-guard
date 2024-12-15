package server_errors

import (
	app_errors "eventsguard/internal/app/errors"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func AppErrorToHumaError(appError app_errors.AppError) error {
	var ret error

	log.Printf("appError.Code: %d", appError.Code)
	log.Printf("appError.Message: %s", appError.Message)

	switch appError.Code {
	case http.StatusNotFound:
		ret = huma.Error404NotFound(appError.Message)
	case http.StatusBadRequest:
		ret = huma.Error400BadRequest(appError.Message)
	case http.StatusUnauthorized:
		ret = huma.Error401Unauthorized(appError.Message)
	case http.StatusForbidden:
		ret = huma.Error403Forbidden(appError.Message)
	case http.StatusNotAcceptable:
		ret = huma.Error406NotAcceptable(appError.Message)
	case http.StatusConflict:
		ret = huma.Error409Conflict(appError.Message)
	case http.StatusPreconditionFailed:
		ret = huma.Error412PreconditionFailed(appError.Message)
	case http.StatusUnprocessableEntity:
		ret = huma.Error422UnprocessableEntity(appError.Message)
	default:
		ret = huma.Error500InternalServerError(appError.Message)
	}

	return ret

}

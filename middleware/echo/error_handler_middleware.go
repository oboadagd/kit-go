package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oboadagd/kit-go/middleware/responses"
	"github.com/pkg/errors"
)

type (
	ErrorHandlerMiddlewareInterface interface {
		HandlerError(next echo.HandlerFunc) echo.HandlerFunc
	}

	ErrorHandlerMiddleware struct {
	}

	// Map defines a generic map of type `map[string]interface{}`.
	Map map[string]interface{}
)

func NewErrorHandlerMiddleware() ErrorHandlerMiddlewareInterface {
	return &ErrorHandlerMiddleware{}
}

func (h *ErrorHandlerMiddleware) HandlerError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			var er error
			var erCode string
			var status int
			// Build the error response
			switch act := errors.Cause(err).(type) {
			case *responses.GenericHttpError:
				er = act.ErrorMsg
				erCode = act.ErrorCode
				status = act.Status
			default:
				er = errors.New(http.StatusText(http.StatusInternalServerError))
				erCode = "server_error"
				status = http.StatusInternalServerError
			}
			return c.JSON(status, Map{"error_code": erCode, "error_message": er.Error()})
		}
		return nil
	}
}

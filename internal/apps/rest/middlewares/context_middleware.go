package middlewares

import (
	"context"
	"go-boilerplate/internal/constants"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ContextTimeoutMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), constants.TimeoutDuration)
		defer cancel()

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func GenerateRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check whether request id already exist on header
		listRequestIDKey := []string{constants.RequestIDKey}
		requestID := ""
		for _, k := range listRequestIDKey {
			if headerRequestID := c.Request().Header.Get(k); headerRequestID != "" {
				requestID = headerRequestID
				break
			}
		}

		// Generate new request id if empty
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request id on request & response header
		c.Request().Header.Set(constants.RequestIDKey, requestID)
		c.Response().Header().Set(constants.RequestIDKey, requestID)

		// Set request id to context
		ctx := c.Request().Context()
		c.SetRequest(c.Request().Clone(context.WithValue(ctx, constants.RequestIDKey, requestID)))

		return next(c)
	}
}

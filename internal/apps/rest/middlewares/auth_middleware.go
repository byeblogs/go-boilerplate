package middlewares

import (
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/pkg/helpers"
	"go-boilerplate/internal/pkg/responses"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeaderList, ok := c.Request().Header["Authorization"]
		if !ok || len(authHeaderList) == 0 {
			return responses.NewError().
				WithError(constants.ErrNotLoggedIn).
				WithCode(http.StatusUnauthorized).
				WithMessage("You don't have the permission.").
				SendErrorResponse(c)
		}

		authHeader := authHeaderList[0]
		bearerPrefix := "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return responses.NewError().
				WithError(constants.ErrNotLoggedIn).
				WithCode(http.StatusUnauthorized).
				WithMessage("Invalid authorization header.").
				SendErrorResponse(c)
		}

		token := strings.Replace(authHeader, bearerPrefix, "", 1)
		claims, err := helpers.ParseAndValidateJWT(token)
		if err != nil {
			return responses.NewError().
				WithError(err).
				WithCode(http.StatusUnauthorized).
				WithMessage("Invalid authorization header.").
				SendErrorResponse(c)
		}

		c.Set(constants.AuthClaimsKey, claims)
		c.Set(constants.AccessToken, token)

		return next(c)
	}
}

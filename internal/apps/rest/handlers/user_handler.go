package handlers

import (
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/dtos"
	"go-boilerplate/internal/pkg/helpers"
	"go-boilerplate/internal/pkg/responses"
	"go-boilerplate/internal/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

type UserHandler interface {
	CreateUser(c echo.Context) (err error)
	GetUserByID(c echo.Context) (err error)
	UpdateUser(c echo.Context) (err error)
	DeleteUser(c echo.Context) (err error)
}

type UserHandlerImpl struct {
	usecase *usecases.Usecase
}

func NewUserHandler(ioc di.Container) UserHandler {
	return &UserHandlerImpl{
		usecase: ioc.Get(constants.Usecase).(*usecases.Usecase),
	}
}

func (t *UserHandlerImpl) CreateUser(c echo.Context) (err error) {
	var (
		ctx    = c.Request().Context()
		params dtos.CreateUserRequest
	)

	if err = c.Bind(&params); err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Failed to bind parameters").
			SendErrorResponse(c)
		return
	}

	err = t.usecase.User.CreateUser(ctx, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusCreated).
		WithMessage("Successfully created a new user").
		Send(c)
}

func (t *UserHandlerImpl) GetUserByID(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		params dtos.UserIDParams
		err    error
	)

	if err = c.Bind(&params); err != nil {
		return responses.NewError().
			WithCode(http.StatusBadRequest).
			WithError(err).
			WithMessage("Failed to bind parameters").
			SendErrorResponse(c)
	}

	claims, err := helpers.GetAuthClaims(c)
	if err != nil {
		return responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Failed to get auth claims")
	}

	data, err := t.usecase.User.GetUserByID(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully retrieved a user").
		WithData(data).
		Send(c)
}

func (t *UserHandlerImpl) UpdateUser(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		params dtos.UpdateUserParams
		err    error
	)

	if err = c.Bind(&params); err != nil {
		return responses.NewError().
			WithCode(http.StatusBadRequest).
			WithError(err).
			WithMessage("Failed to bind parameters").
			SendErrorResponse(c)
	}

	claims, err := helpers.GetAuthClaims(c)
	if err != nil {
		return responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Failed to get auth claims")
	}

	err = t.usecase.User.UpdateUser(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully updated a user").
		Send(c)
}

func (t *UserHandlerImpl) DeleteUser(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		params dtos.UserIDParams
		err    error
	)

	if err = c.Bind(&params); err != nil {
		return responses.NewError().
			WithCode(http.StatusBadRequest).
			WithError(err).
			WithMessage("Failed to bind parameters").
			SendErrorResponse(c)
	}

	claims, err := helpers.GetAuthClaims(c)
	if err != nil {
		return responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Failed to get auth claims")
	}

	err = t.usecase.User.DeleteUser(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully deleted a user").
		Send(c)
}

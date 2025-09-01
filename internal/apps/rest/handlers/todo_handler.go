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

type TodoHandler interface {
	CreateTodo(c echo.Context) (err error)
	GetTodoByID(c echo.Context) (err error)
	GetTodos(c echo.Context) (err error)
	UpdateTodo(c echo.Context) (err error)
	DeleteTodo(c echo.Context) (err error)
}

type TodoHandlerImpl struct {
	usecase *usecases.Usecase
}

func NewTodoHandler(ioc di.Container) TodoHandler {
	return &TodoHandlerImpl{
		usecase: ioc.Get(constants.Usecase).(*usecases.Usecase),
	}
}

func (t *TodoHandlerImpl) CreateTodo(c echo.Context) (err error) {
	var (
		ctx    = c.Request().Context()
		params dtos.CreateTodoRequest
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

	err = t.usecase.Todo.CreateTodo(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusCreated).
		WithMessage("Successfully created a new todo").
		Send(c)
}

func (t *TodoHandlerImpl) GetTodoByID(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		params dtos.TodoIDParams
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

	data, err := t.usecase.Todo.GetTodoByID(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully retrieved a todo").
		WithData(data).
		Send(c)
}

func (t *TodoHandlerImpl) GetTodos(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	claims, err := helpers.GetAuthClaims(c)
	if err != nil {
		return responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Failed to get auth claims")
	}

	data, err := t.usecase.Todo.GetTodos(ctx, claims)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully retrieved todos").
		WithData(data).
		Send(c)
}

func (t *TodoHandlerImpl) UpdateTodo(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		params dtos.UpdateTodoParams
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

	err = t.usecase.Todo.UpdateTodo(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully updated a todo").
		Send(c)
}

func (t *TodoHandlerImpl) DeleteTodo(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		params dtos.TodoIDParams
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

	err = t.usecase.Todo.DeleteTodo(ctx, claims, params)
	return responses.New().
		WithError(err).
		WithSuccessCode(http.StatusOK).
		WithMessage("Successfully deleted a todo").
		Send(c)
}

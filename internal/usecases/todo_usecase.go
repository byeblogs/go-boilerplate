package usecases

import (
	"context"
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/dtos"
	"go-boilerplate/internal/models"
	"go-boilerplate/internal/pkg/responses"
	"go-boilerplate/internal/pkg/utils"
	"go-boilerplate/internal/repositories"
	"net/http"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

type TodoUseCase interface {
	CreateTodo(ctx context.Context, claims dtos.AuthClaims, params dtos.CreateTodoRequest) (err error)
	GetTodoByID(ctx context.Context, claims dtos.AuthClaims, params dtos.TodoIDParams) (data dtos.GetTodoByIDResponse, err error)
	GetTodos(ctx context.Context, claims dtos.AuthClaims) (data dtos.GetTodosResponse, err error)
	UpdateTodo(ctx context.Context, claims dtos.AuthClaims, params dtos.UpdateTodoParams) (err error)
	DeleteTodo(ctx context.Context, claims dtos.AuthClaims, params dtos.TodoIDParams) (err error)
}

type TodoUseCaseImpl struct {
	repository	*repositories.Repository
	util		*utils.Util
}

func NewTodoUseCase(ioc di.Container) *TodoUseCaseImpl {
	return &TodoUseCaseImpl{
		repository: ioc.Get(constants.Repository).(*repositories.Repository),
		util:       ioc.Get(constants.Util).(*utils.Util),
	}
}

func (s *TodoUseCaseImpl) CreateTodo(ctx context.Context, claims dtos.AuthClaims, params dtos.CreateTodoRequest) (err error) {
	user, err := s.repository.User.GetUserByID(ctx, claims.UserID)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while retrieving user by id from database")
		return
	}

	if user == nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Cannot find user with the given id")
		return
	}

	currentTime, err := s.util.Date.GetTimeNowJakarta()
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while getting current time")
		return
	}

	newTodo := models.Todo{
		Title:   params.Title,
		Content: params.Content,
		UserID:  user.ID,
		Model: gorm.Model{
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		},
	}

	err = s.repository.Todo.CreateTodo(ctx, newTodo)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while creating todo into database")
		return
	}

	return
}

func (s *TodoUseCaseImpl) GetTodoByID(ctx context.Context, claims dtos.AuthClaims, params dtos.TodoIDParams) (data dtos.GetTodoByIDResponse, err error) {
	todo, err := s.repository.Todo.GetTodoByID(ctx, params.ID)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while retrieving todo by id from database")
		return
	}

	if todo == nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Cannot find todo with the given id")
		return
	}

	if todo.UserID != claims.UserID {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusUnauthorized).
			WithMessage("You are not authorized to view this todo")
		return
	}

	data.Todo = *todo
	return
}

func (s *TodoUseCaseImpl) GetTodos(ctx context.Context, claims dtos.AuthClaims) (data dtos.GetTodosResponse, err error) {
	todos, err := s.repository.Todo.GetTodosByUserID(ctx, claims.UserID)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while retrieving todos from database")
		return
	}

	data.Todos = todos
	return
}

func (s *TodoUseCaseImpl) UpdateTodo(ctx context.Context, claims dtos.AuthClaims, params dtos.UpdateTodoParams) (err error) {
	todo, err := s.repository.Todo.GetTodoByID(ctx, params.ID)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while retrieving todo by id from database")
		return
	}

	if todo == nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Cannot find todo with the given id")
		return
	}

	if todo.UserID != claims.UserID {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusUnauthorized).
			WithMessage("You are not authorized to update this todo")
		return
	}

	todo.Title = params.Title
	todo.Content = params.Content

	currentTime, err := s.util.Date.GetTimeNowJakarta()
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while getting current time")
		return
	}
	todo.Model.UpdatedAt = currentTime

	err = s.repository.Todo.UpdateTodo(ctx, *todo)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while updating todo into database")
		return
	}

	return
}

func (s *TodoUseCaseImpl) DeleteTodo(ctx context.Context, claims dtos.AuthClaims, params dtos.TodoIDParams) (err error) {
	todo, err := s.repository.Todo.GetTodoByID(ctx, params.ID)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Cannot find todo with the given id")
		return
	}

	if todo == nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Cannot find todo with the given id")
		return
	}

	if todo.UserID != claims.UserID {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusUnauthorized).
			WithMessage("You are not authorized to delete this todo")
		return
	}

	err = s.repository.Todo.DeleteTodo(ctx, *todo)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while deleting todo from database")
		return
	}

	return
}

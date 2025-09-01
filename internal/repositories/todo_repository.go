package repositories

import (
	"context"
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/models"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, todo models.Todo) (err error)
	GetTodoByID(ctx context.Context, todoID uint) (todo *models.Todo, err error)
	GetTodosByUserID(ctx context.Context, userID uint) (todos []models.Todo, err error)
	UpdateTodo(ctx context.Context, todo models.Todo) (err error)
	DeleteTodo(ctx context.Context, todo models.Todo) (err error)
}

type TodoRepositoryImpl struct {
	db *gorm.DB
}

func NewTodoRepository(ioc di.Container) *TodoRepositoryImpl {
	return &TodoRepositoryImpl{
		db: ioc.Get(constants.Postgres).(*gorm.DB),
	}
}

func (r *TodoRepositoryImpl) CreateTodo(ctx context.Context, todo models.Todo) error {
	err := r.db.Create(&todo).WithContext(ctx).Error
	return err
}

func (r *TodoRepositoryImpl) GetTodoByID(ctx context.Context, todoID uint) (todo *models.Todo, err error) {
	err = r.db.Where("id = ?", todoID).Find(&todo).Limit(1).WithContext(ctx).Error
	if todo.ID == 0 {
		return nil, nil
	}
	return
}

func (r *TodoRepositoryImpl) GetTodosByUserID(ctx context.Context, userID uint) (todos []models.Todo, err error) {
	err = r.db.Where("user_id = ?", userID).Find(&todos).WithContext(ctx).Error
	return
}

func (r *TodoRepositoryImpl) UpdateTodo(ctx context.Context, todo models.Todo) error {
	err := r.db.Save(&todo).WithContext(ctx).Error
	return err
}

func (r *TodoRepositoryImpl) DeleteTodo(ctx context.Context, todo models.Todo) error {
	err := r.db.Delete(&todo).WithContext(ctx).Error
	return err
}

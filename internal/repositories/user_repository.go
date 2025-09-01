package repositories

import (
	"context"
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/models"

	"github.com/sarulabs/di"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (err error)
	GetUserByID(ctx context.Context, userID uint) (user *models.User, err error)
	GetUserByUsername(ctx context.Context, username string) (user *models.User, err error)
	UpdateUser(ctx context.Context, user models.User) (err error)
	DeleteUser(ctx context.Context, user models.User) (err error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(ioc di.Container) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: ioc.Get(constants.Postgres).(*gorm.DB),
	}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user models.User) (err error) {
	err = r.db.Create(&user).WithContext(ctx).Error
	return
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID uint) (user *models.User, err error) {
	err = r.db.Where("id = ?", userID).Find(&user).Limit(1).WithContext(ctx).Error
	if user.ID == 0 {
		return nil, nil
	}
	return
}

func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (user *models.User, err error) {
	err = r.db.Where("username = ?", username).Find(&user).Limit(1).WithContext(ctx).Error
	if user.ID == 0 {
		return nil, nil
	}
	return
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user models.User) (err error) {
	err = r.db.Save(&user).WithContext(ctx).Error
	return
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, user models.User) (err error) {
	err = r.db.Delete(&user).WithContext(ctx).Error
	return
}

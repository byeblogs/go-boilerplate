package mock_repositories

import (
	"go-boilerplate/internal/mocks"
	"go-boilerplate/internal/repositories"

	"github.com/golang/mock/gomock"
)

var (
	UserRepository    *mocks.MockUserRepository
)

func NewMockRepository(ctrl *gomock.Controller) *repositories.Repository {
	UserRepository = mocks.NewMockUserRepository(ctrl)

	return &repositories.Repository{
		User: UserRepository,
	}
}
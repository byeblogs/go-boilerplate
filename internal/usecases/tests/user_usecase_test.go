package tests

import (
	"context"
	"go-boilerplate/internal/dtos"
	"go-boilerplate/internal/mock"
	mock_repositories "go-boilerplate/internal/mock/repositories"
	"go-boilerplate/internal/models"
	"go-boilerplate/internal/usecases"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserServiceSuite struct {
	suite.Suite
	ctx         context.Context
	ctrl        *gomock.Controller
	userService usecases.UserUseCase
}

func (s *UserServiceSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	module := mock.ModuleMock(s.ctrl)
	s.userService = usecases.NewUserUseCase(module)
	s.ctx = context.Background()
}

func (s *UserServiceSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestUserUseCase(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) TestCreateUser() {
	s.Run("Success", func() {
		s.SetupTest()
		mock_repositories.UserRepository.EXPECT().GetUserByUsername(gomock.Any(), "rahmat").Return(nil, nil)
		mock_repositories.UserRepository.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)

		err := s.userService.CreateUser(s.ctx, dtos.CreateUserRequest{Username: "rahmat", Password: "rahmat"})

		s.Equal(nil, err)
	})
}

func (s *UserServiceSuite) TestGetUserByID() {
	s.Run("Success With No Error", func() {
		s.SetupTest()
		mock_repositories.UserRepository.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}}, nil)

		data, err := s.userService.GetUserByID(s.ctx, dtos.AuthClaims{UserID: 1}, dtos.UserIDParams{ID: 1})

		s.Equal(nil, err)
		s.Equal(dtos.GetUserByIDResponse{User: models.User{Model: gorm.Model{ID: 1}}}, data)
	})
}

func (s *UserServiceSuite) TestUpdateUser() {
	s.Run("Success With No Error", func() {
		s.SetupTest()
		mock_repositories.UserRepository.EXPECT().GetUserByID(gomock.Any(), uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}}, nil)
		mock_repositories.UserRepository.EXPECT().UpdateUser(gomock.Any(), models.User{Model: gorm.Model{ID: 1}, FullName: "John Wick", PhoneNumber: "08981297512"})

		err := s.userService.UpdateUser(s.ctx, dtos.AuthClaims{UserID: 1}, dtos.UpdateUserParams{ID: 1, FullName: "John Wick", PhoneNumber: "08981297512"})

		s.Equal(nil, err)
	})
}

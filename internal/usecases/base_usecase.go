package usecases

import "github.com/sarulabs/di"

type Usecase struct {
	User	UserUseCase
	Auth	AuthUseCase
	Todo	TodoUseCase
}

func NewUsecase(ioc di.Container) *Usecase {
	return &Usecase{
		User: NewUserUseCase(ioc),
		Auth: NewAuthUseCase(ioc),
		Todo: NewTodoUseCase(ioc),
	}
}
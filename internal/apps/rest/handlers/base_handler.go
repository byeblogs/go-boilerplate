package handlers

import "github.com/sarulabs/di"

type Controller struct {	
	User	UserHandler
	Auth	AuthHandler
	Todo	TodoHandler
}

func NewController(ioc di.Container) *Controller {
	return &Controller{
		User: NewUserHandler(ioc),
		Auth: NewAuthHandler(ioc),
		Todo: NewTodoHandler(ioc),
	}
}
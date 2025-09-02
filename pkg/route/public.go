package route

import (
	"github.com/byeblogs/go-boilerplate/app/controller"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public route.
func PublicRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/api/v1")

	route.Post("/token/new", controller.GetNewAccessToken)
	route.Get("/books", controller.GetBooks)
	route.Get("/books/:id", controller.GetBook)

}

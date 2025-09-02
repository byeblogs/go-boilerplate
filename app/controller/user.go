package controller

import (
	"strconv"

	"github.com/byeblogs/go-boilerplate/app/dto"
	"github.com/byeblogs/go-boilerplate/app/model"
	repo "github.com/byeblogs/go-boilerplate/app/repository"
	"github.com/byeblogs/go-boilerplate/pkg/validator"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetUsers func gets all exists user.
// @Description Get all exists user.
// @Summary get all exists user
// @Tags User
// @Accept json
// @Produce json
// @Param page query integer false "Page no"
// @Param page_size query integer false "records per page"
// @Success 200 {object} dto.User "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Security ApiKeyAuth
// @Router /v1/users [get]
func GetUsers(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)
	userRepo := repo.NewUserRepo(database.GetDB())
	users, err := userRepo.All(pageSize, uint(pageSize*(pageNo-1)))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "users were not found",
		})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(users),
		"users":     dto.ToUsers(users),
	})
}

// GetUser func gets a user.
// @Description a user.
// @Summary get a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.User "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Security ApiKeyAuth
// @Router /v1/users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	ID, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	userRepo := repo.NewUserRepo(database.GetDB())
	user, err := userRepo.Get(int(ID))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user were not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(user),
	})
}

// CreateUser func for creates a new user.
// @Description Create a new user.
// @Summary create a new user
// @Tags User
// @Accept json
// @Produce json
// @Param createuser body model.CreateUser true "Create new user"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Success 200 {object} dto.User "Ok" status "Ok"
// @Security ApiKeyAuth
// @Router /v1/users [post]
func CreateUser(c *fiber.Ctx) error {
	// Create new Book struct
	user := &model.CreateUser{}

	if err := c.BodyParser(user); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := validator.NewValidator()
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	// check user already exists
	exists, err := userRepo.Exists(user.UserName, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"msg": "user with this username or email already exists",
		})
	}

	user.Password, _ = GeneratePasswordHash([]byte(user.Password))
	if err := userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dbUser, err := userRepo.GetByUsername(user.UserName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(dbUser),
	})
}

// UpdateUser func update a user.
// @Description first_name, last_name, is_active, is_admin only
// @Summary update a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param updateuser body model.UpdateUser true "Update a user"
// @Success 200 {object} dto.User "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Security ApiKeyAuth
// @Router /v1/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	ID64, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	ID := int(ID64)
	userRepo := repo.NewUserRepo(database.GetDB())
	_, err = userRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user were not found",
		})
	}

	user := &model.UpdateUser{}
	if err := c.BodyParser(user); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := validator.NewValidator()
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	if err := userRepo.Update(ID, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dbUser, err := userRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": dto.ToUser(dbUser),
	})
}

// DeleteUser func delete a user.
// @Description delete user
// @Summary delete a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} interface{} "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Security ApiKeyAuth
// @Router /v1/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	ID64, err := strconv.ParseInt(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	ID := int(ID64)
	userRepo := repo.NewUserRepo(database.GetDB())
	_, err = userRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "user were not found",
		})
	}

	err = userRepo.Delete(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{})
}

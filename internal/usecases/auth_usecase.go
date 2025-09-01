package usecases

import (
	"context"
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/dtos"
	"go-boilerplate/internal/pkg/helpers"
	"go-boilerplate/internal/pkg/responses"
	"go-boilerplate/internal/pkg/utils"
	"go-boilerplate/internal/repositories"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sarulabs/di"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Login(ctx context.Context, params dtos.LoginRequest) (data dtos.LoginResponse, err error)
	Logout(ctx context.Context, authClaims dtos.AuthClaims) (err error)
}

type AuthUseCaseImpl struct {
	repository *repositories.Repository
	util       *utils.Util
}

func NewAuthUseCase(ioc di.Container) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		repository: ioc.Get(constants.Repository).(*repositories.Repository),
		util:       ioc.Get(constants.Util).(*utils.Util),
	}
}

func (s *AuthUseCaseImpl) Login(ctx context.Context, params dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	user, err := s.repository.User.GetUserByUsername(ctx, params.Username)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while retrieving user by username from database")
		return
	}

	if user == nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Cannot find user with the given username")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusBadRequest).
			WithMessage("Incorrect password")
		return
	}

	tokenExpireDuration := (time.Hour * 24)
	currentTime, err := s.util.Date.GetTimeNowJakarta()
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while getting current time")
		return
	}

	token, err := helpers.GenerateJWTString(dtos.AuthClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(tokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
		},
	})
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while generating JWT")
		return
	}

	res.AccessToken = token
	return
}

func (s *AuthUseCaseImpl) Logout(ctx context.Context, authClaims dtos.AuthClaims) (err error) {
	user, err := s.repository.User.GetUserByID(ctx, authClaims.UserID)
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

	err = s.repository.User.UpdateUser(ctx, *user)
	if err != nil {
		err = responses.NewError().
			WithError(err).
			WithCode(http.StatusInternalServerError).
			WithMessage("Error while updating users access token into database")
		return
	}

	return
}

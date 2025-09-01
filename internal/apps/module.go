package apps

import (
	"go-boilerplate/internal/apps/event/subscribers"
	"go-boilerplate/internal/apps/rest/handlers"
	"go-boilerplate/internal/constants"
	"go-boilerplate/internal/pkg/databases"
	"go-boilerplate/internal/pkg/utils"
	"go-boilerplate/internal/repositories"
	"go-boilerplate/internal/usecases"

	"github.com/sarulabs/di"
)

func NewIOC() di.Container {
	builder, _ := di.NewBuilder()
	_ = builder.Add(
		di.Def{
			Name: constants.Postgres,
			Build: func(ctn di.Container) (interface{}, error) {
				db, err := databases.NewPostgresClient()
				return db, err
			},
		},
		di.Def{
			Name: constants.Redis,
			Build: func(ctn di.Container) (interface{}, error) {
				rdb, err := databases.NewRedisClient()
				return rdb, err
			},
		},
		di.Def{
			Name: constants.Controller,
			Build: func(ctn di.Container) (interface{}, error) {
				return handlers.NewController(builder.Build()), nil
			},
		},
		di.Def{
			Name: constants.Usecase,
			Build: func(ctn di.Container) (interface{}, error) {
				return usecases.NewUsecase(builder.Build()), nil
			},
		},
		di.Def{
			Name: constants.Repository,
			Build: func(ctn di.Container) (interface{}, error) {
				return repositories.NewRepository(builder.Build()), nil
			},
		},
		di.Def{
			Name: constants.Util,
			Build: func(ctn di.Container) (interface{}, error) {
				return utils.NewUtil(), nil
			},
		},
		di.Def{
			Name: constants.Subscriber,
			Build: func(ctn di.Container) (interface{}, error) {
				return subscribers.NewSubscriber(builder.Build()), nil
			},
		},
	)
	return builder.Build()
}

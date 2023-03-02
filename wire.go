//go:build wireinject
// +build wireinject

package api

import (
	"api/auth"
	"api/auth/cognito"
	"api/core"
	"api/image"
	"api/image/resize"
	"api/storage"
	"api/storage/postgresql"
	"github.com/google/wire"
	"github.com/rs/zerolog"
)

var DatabaseSet = wire.NewSet(
	postgresql.NewDatabase,
	postgresql.NewImageRepository,
	postgresql.NewUserRepo,
	wire.Bind(new(storage.Storage), new(*postgresql.Database)),
	wire.Bind(new(storage.ImagesRepository), new(*postgresql.ImageRepo)),
	wire.Bind(new(storage.UserRepository), new(*postgresql.UserRepo)),
)

func InitializeApp(logger *zerolog.Logger) (*core.App, error) {
	wire.Build(
		core.NewConfigFromEnv,
		DatabaseSet,
		resize.NewClient,
		wire.Bind(new(image.Resizer), new(*resize.Client)),
		cognito.NewCognitoAuthService,
		wire.Bind(new(auth.Authenticator), new(*cognito.AuthService)),
		core.NewImagesService,
		core.NewApp,
	)

	return new(core.App), nil
}

func InitializeAppForTesting(logger *zerolog.Logger) (*core.App, error) {
	wire.Build(
		core.NewConfigFromEnv,
		DatabaseSet,
		resize.NewClient,
		wire.Bind(new(image.Resizer), new(*resize.Client)),
		cognito.NewCognitoAuthService,
		wire.Bind(new(auth.Authenticator), new(*cognito.AuthService)),
		core.NewImagesService,
		core.NewApp,
	)

	return new(core.App), nil
}

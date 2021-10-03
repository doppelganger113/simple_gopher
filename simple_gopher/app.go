package simple_gopher

import (
	"context"
	"simple_gopher/auth"
	"simple_gopher/auth/cognito"
	"simple_gopher/image/image_resize_api"
	"simple_gopher/storage"
	"simple_gopher/storage/postgresql"
)

type App struct {
	Config        Config
	ImagesService ImagesService
	Auth          auth.Authenticator
	storage       storage.Storage
}

func NewApp(
	config Config,
	storage storage.Storage,
	auth auth.Authenticator,
	imagesService ImagesService,
) *App {
	return &App{
		Config:        config,
		storage:       storage,
		Auth:          auth,
		ImagesService: imagesService,
	}
}

func (a *App) Init(initCtx context.Context, ctx context.Context) error {
	err := a.storage.Connect(initCtx, a.Config.DatabaseUrl)
	if err != nil {
		return err
	}
	err = a.Auth.FetchAndSetKeySet(initCtx)
	if err != nil {
		return err
	}

	if !a.Config.SqsPostAuthConsumerDisabled {
		a.Auth.StartConsumingPostAuthAsync(ctx)
	}

	return nil
}

func (a *App) Shutdown(_ context.Context) error {
	a.storage.Close()
	if !a.Config.SqsPostAuthConsumerDisabled {
		return a.Auth.Shutdown()
	}
	return nil
}

func CreateApp(config Config) (*App, error) {
	authConfig := NewAuthConfig(config)

	// Storage
	db := postgresql.NewDatabase()
	imageRepository := postgresql.NewImageRepository(db)
	userRepository := postgresql.NewUserRepo(db)

	// Services
	imageResizeConfig := image_resize_api.Config{ImagesApiDomain: config.ImagesApiDomain}
	resizeApi := image_resize_api.NewResizeApi(imageResizeConfig)
	authenticator := cognito.NewCognitoAuthService(authConfig, userRepository)
	imagesService := NewImagesService(resizeApi, imageRepository, authenticator)

	// App
	app := NewApp(
		config, db, authenticator, imagesService,
	)

	return app, nil
}

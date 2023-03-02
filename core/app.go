package core

import (
	"api/auth"
	"api/storage"
	"context"
	"fmt"
)

type App struct {
	Config        Config
	ImagesService *ImagesService
	Auth          auth.Authenticator
	storage       storage.Storage
}

func NewApp(
	config Config,
	storage storage.Storage,
	auth auth.Authenticator,
	imagesService *ImagesService,
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
		return fmt.Errorf("failed connecting to storage: %w", err)
	}
	err = a.Auth.FetchAndSetKeySet(initCtx)
	if err != nil {
		return fmt.Errorf("failed fetching and setting authentication key set: %w", err)
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

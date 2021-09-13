package storage

import (
	"context"
)

type ImagesRepository interface {
	Get(ctx context.Context, limit, offset int, order Order) (ImageList, error)
	GetOne(ctx context.Context, imageId string) (Image, error)
	GetOneByName(ctx context.Context, name string) (Image, error)
	DoesImageExist(ctx context.Context, name string) (bool, error)
	Create(ctx context.Context, image Image) (Image, error)
	SetNameById(ctx context.Context, imageId, newName string) (Image, error)
	UpdateOne(ctx context.Context, updates Image) error
	DeleteOne(ctx context.Context, imageId string) error
}

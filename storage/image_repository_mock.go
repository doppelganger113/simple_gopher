package storage

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ImageRepoMock struct {
	mock.Mock
}

func (repo *ImageRepoMock) Get(
	ctx context.Context, page, size int, order Order,
) (ImageList, error) {
	args := repo.Called(ctx, page, size, order)

	return args.Get(0).(ImageList), args.Error(1)
}

func (repo *ImageRepoMock) GetOne(ctx context.Context, imageId string) (*Image, error) {
	args := repo.Called(ctx, imageId)

	return args.Get(0).(*Image), args.Error(1)
}

func (repo *ImageRepoMock) GetOneByName(ctx context.Context, name string) (Image, error) {
	args := repo.Called(ctx, name)

	return args.Get(0).(Image), args.Error(1)
}

func (repo *ImageRepoMock) DoesImageExist(ctx context.Context, name string) (bool, error) {
	args := repo.Called(ctx, name)
	return args.Bool(0), args.Error(1)
}

func (repo *ImageRepoMock) Create(ctx context.Context, newImage Image) (*Image, error) {
	args := repo.Called(ctx, newImage)

	return args.Get(0).(*Image), args.Error(1)
}

func (repo *ImageRepoMock) SetNameById(ctx context.Context, imageId, newName string) (Image, error) {
	args := repo.Called(ctx, imageId, newName)

	return args.Get(0).(Image), args.Error(1)
}

func (repo *ImageRepoMock) UpdateOne(ctx context.Context, updates Image) error {
	args := repo.Called(ctx, updates)

	return args.Error(0)
}

func (repo *ImageRepoMock) DeleteOne(ctx context.Context, imageId string) error {
	args := repo.Called(ctx, imageId)

	return args.Error(0)
}

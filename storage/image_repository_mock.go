package storage

import (
	"context"
)

type ImageRepoMock struct {
}

func (repo ImageRepoMock) Get(
	_ context.Context, _, _ int, _ Order,
) (ImageList, error) {
	images := ImageList{
		{
			Id:       "3c47d736-6c4e-4a1c-a04b-3744cc30b263",
			Name:     "my-image-1",
			Format:   "jpg",
			Original: "images/my-unique-image.png",
			Domain:   "https://random.cloudfront.net",
			Path:     "images",
			Sizes: ImageSizes{
				Original: Dimensions{
					Width:  688,
					Height: 516,
				},
				Xs: &Dimensions{
					Width:  100,
					Height: 75,
				},
				S: &Dimensions{
					Width:  300,
					Height: 225,
				},
				M: &Dimensions{
					Width:  500,
					Height: 375,
				},
			},
			CreatedAt: nil,
			UpdatedAt: nil,
		},
	}

	return images, nil
}

func (repo ImageRepoMock) GetOne(_ context.Context, _ string) (Image, error) {
	return Image{}, nil
}

func (repo ImageRepoMock) GetOneByName(_ context.Context, _ string) (Image, error) {
	return Image{}, nil
}

func (repo ImageRepoMock) DoesImageExist(_ context.Context, _ string) (bool, error) {
	return false, nil
}

func (repo ImageRepoMock) Create(_ context.Context, _ Image) (Image, error) {
	return Image{}, nil
}

func (repo ImageRepoMock) SetNameById(_ context.Context, _, _ string) (Image, error) {
	return Image{}, nil
}

func (repo ImageRepoMock) UpdateOne(_ context.Context, _ Image) error {
	return nil
}

func (repo ImageRepoMock) DeleteOne(_ context.Context, _ string) error {
	return nil
}

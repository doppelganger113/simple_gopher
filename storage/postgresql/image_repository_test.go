package postgresql

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"simple_gopher/image"
	"simple_gopher/storage"
	"simple_gopher/test"
	"testing"
)

func setupImageRepo(ctx context.Context) (*ImageRepo, error) {
	db, err := setupDb(ctx)
	if err != nil {
		return nil, err
	}

	repo := NewImageRepository(db)

	return repo, nil
}

func cleanImageRepo(t *testing.T, repo *ImageRepo) {
	defer repo.database.Close()

	_, err := repo.DeleteAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func insertDummyData(repo *ImageRepo, userRepo *UserRepo) error {
	ctx := context.Background()

	createdUser, err := userRepo.Create(ctx, storage.UserCreationDto{
		Email:       "john@gmail.com",
		Role:        "",
		CogUsername: "",
		CogSub:      "",
		CogName:     "",
		Disabled:    false,
	})
	if err != nil {
		return err
	}

	imageList := storage.ImageList{
		{
			Id:       "",
			Name:     "testing-image-one",
			Format:   "png",
			Original: "images/testing-image-one.png",
			Domain:   "https://whatever.com",
			Path:     "images",
			AuthorId: createdUser.Id,
			Sizes: storage.ImageSizes{
				Original: storage.Dimensions{},
			},
		},
		{
			Id:       "",
			Name:     "testing-image-two",
			Format:   "jpg",
			Original: "images/testing-image-two.jpg",
			Domain:   "https://whatever.com",
			Path:     "images",
			AuthorId: createdUser.Id,
			Sizes: storage.ImageSizes{
				Original: storage.Dimensions{Width: 300, Height: 400},
				Xs:       &storage.Dimensions{Width: 100, Height: 150},
				S:        &storage.Dimensions{Width: 200, Height: 250},
			},
		},
	}

	_, err = repo.InsertMany(ctx, imageList)
	return err
}

func TestImageRepository_Get(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupImageRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	userRepo, err := setupUserRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUserRepo(userRepo)
	defer cleanImageRepo(t, repo)

	err = insertDummyData(repo, userRepo)
	if err != nil {
		t.Error(fmt.Errorf("error inserting images %w", err))
	}

	imageList, err := repo.Get(context.Background(), 10, 0, storage.OrderDescending)
	if err != nil {
		t.Fatal(err)
	}

	if imageList == nil {
		t.Fatal("Got nil images")
	}
	if len(imageList) != 2 {
		t.Fatalf("Expected 2 images, got %d", len(imageList))
	}

	expected := storage.ImageList{
		{
			Id:       "",
			Name:     "testing-image-one",
			Format:   "png",
			Original: "images/testing-image-one.png",
			Domain:   "https://whatever.com",
			Path:     "images",
			AuthorId: "user-1",
			Sizes: storage.ImageSizes{
				Original: storage.Dimensions{},
			},
		},
		{
			Id:       "",
			Name:     "testing-image-two",
			Format:   "jpg",
			Original: "images/testing-image-two.jpg",
			Domain:   "https://whatever.com",
			Path:     "images",
			AuthorId: "user-1",
			Sizes: storage.ImageSizes{
				Original: storage.Dimensions{Width: 300, Height: 400},
				Xs:       &storage.Dimensions{Width: 100, Height: 150},
				S:        &storage.Dimensions{Width: 200, Height: 250},
			},
		},
	}
	if imageList.IsEqualTo(expected) == false {
		t.Fatalf(
			"Images are not as expected!\nExpected: \n%s\nGot: \n%s",
			expected.ToString(),
			imageList.ToString(),
		)
	}
}

func TestImageRepo_GetOne(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupImageRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	userRepo, err := setupUserRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUserRepo(userRepo)
	defer cleanImageRepo(t, repo)

	err = insertDummyData(repo, userRepo)
	if err != nil {
		t.Error(fmt.Errorf("error inserting images %w", err))
	}

	imageList, err := repo.Get(context.Background(), 10, 0, storage.OrderDescending)
	if err != nil {
		t.Fatal(err)
	}

	if imageList == nil {
		t.Fatal("Got nil images")
	}

	image, err := repo.GetOne(context.Background(), imageList[1].Id)
	assert.Nil(t, err)
	assert.Equal(t, imageList[1].Id, image.Id)
	assert.Equal(t, imageList[1].Domain, image.Domain)
	assert.Equal(t, imageList[1].Format, image.Format)
	assert.Equal(t, imageList[1].Name, image.Name)
	assert.Equal(t, imageList[1].Sizes, image.Sizes)
	assert.Equal(t, imageList[1].AuthorId, image.AuthorId)
	assert.Equal(t, imageList[1].CreatedAt, image.CreatedAt)
	assert.Equal(t, imageList[1].Path, image.Path)
	assert.Equal(t, imageList[1].Original, image.Original)
}

func TestImageRepository_Create(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupImageRepo(ctx)
	assert.Nil(t, err)

	userRepo, err := setupUserRepo(ctx)
	assert.Nil(t, err)
	defer cleanUserRepo(userRepo)
	defer cleanImageRepo(t, repo)

	newUser, err := userRepo.Create(ctx, storage.UserCreationDto{
		Email:       "john@gmail.com",
		Role:        "",
		CogUsername: "",
		CogSub:      "",
		CogName:     "",
		Disabled:    false,
	})
	assert.Nil(t, err)

	img := storage.Image{
		Name:     "another-new-name",
		Format:   "png",
		Original: "images/another-new-name.png",
		Domain:   "cloud.net",
		Path:     "images",
		AuthorId: newUser.Id,
		Sizes: storage.ImageSizes{
			Original: storage.Dimensions{
				Width:  500,
				Height: 300,
			},
		},
	}

	createdImage, err := repo.Create(context.Background(), img)
	assert.Nil(t, err)

	assert.Equal(t, "images/another-new-name.png", createdImage.Original)
	assert.Equal(t, "another-new-name", createdImage.Name)
	assert.Equal(t, image.PngFormat, createdImage.Format)
	assert.Equal(t, "cloud.net", createdImage.Domain)
	assert.Equal(t, "images", createdImage.Path)
	assert.Equal(t, newUser.Id, createdImage.AuthorId)
	assert.Equal(t, img.Sizes, createdImage.Sizes)
}

func TestImageRepository_DoesImageExist(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupImageRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	userRepo, err := setupUserRepo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanUserRepo(userRepo)
	defer cleanImageRepo(t, repo)

	err = insertDummyData(repo, userRepo)
	if err != nil {
		t.Error(fmt.Errorf("error inserting images %w", err))
	}

	doesExist, err := repo.DoesImageExist(context.Background(), "testing-image-two")
	if err != nil {
		t.Fatal("[DoesImageExist]: ", err)
	}
	if doesExist == false {
		t.Fatal("image should exist")
	}

	doesExist, err = repo.DoesImageExist(context.Background(), "unknown")
	if err != nil {
		t.Fatal("[DoesImageExist unknown]: ", err)
	}
	if doesExist == true {
		t.Fatal("image unknown should not exist")
	}
}

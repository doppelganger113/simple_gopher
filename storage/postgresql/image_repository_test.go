package postgresql

import (
	"api/storage"
	"api/test"
	"context"
	"fmt"
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

	img, err := repo.GetOne(context.Background(), imageList[1].Id)
	if err != nil {
		t.Fatalf("expected empty error, got %v", err)
	}
	secondImage := imageList[1]
	if secondImage.Id != img.Id {
		t.Fatal("failed id assertion")
	}
	if secondImage.Domain != img.Domain {
		t.Fatal("failed Domain assertion")
	}
	if secondImage.Format != img.Format {
		t.Fatal("failed Format assertion")
	}
	if secondImage.Name != img.Name {
		t.Fatal("failed Name assertion")
	}
	if secondImage.Sizes != img.Sizes {
		t.Fatal("failed Sizes assertion")
	}
	if secondImage.AuthorId != img.AuthorId {
		t.Fatal("failed AuthorId assertion")
	}
	if secondImage.CreatedAt != img.CreatedAt {
		t.Fatal("failed CreatedAt assertion")
	}
	if secondImage.Path != img.Path {
		t.Fatal("failed Path assertion")
	}
	if secondImage.Original != img.Original {
		t.Fatal("failed Original assertion")
	}
}

func TestImageRepository_Create(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupImageRepo(ctx)
	if err != nil {
		t.Errorf("failed setting image repo: %v", err)
	}

	userRepo, err := setupUserRepo(ctx)
	if err != nil {
		t.Errorf("failed setting user repo: %v", err)
	}
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
	if err != nil {
		t.Errorf("failed creating user: %v", err)
	}

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
	if err != nil {
		t.Errorf("failed creating image: %v", err)
	}

	if createdImage.Original != "images/another-new-name.png" {
		t.Fatal("failed asserting Original")
	}
	if createdImage.Name != "another-new-name" {
		t.Fatal("failed asserting Name")
	}
	if createdImage.Format != storage.PngFormat {
		t.Fatal("failed asserting Format")
	}
	if createdImage.Domain != "cloud.net" {
		t.Fatal("failed asserting Domain")
	}
	if createdImage.Path != "images" {
		t.Fatal("failed asserting Path")
	}
	if createdImage.AuthorId != newUser.Id {
		t.Fatal("failed asserting Id")
	}
	if createdImage.Sizes != img.Sizes {
		t.Fatal("failed asserting Sizes")
	}
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

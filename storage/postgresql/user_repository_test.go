package postgresql

import (
	"api/storage"
	"api/test"
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

func setupUserRepo(ctx context.Context) (*UserRepo, error) {
	db, err := setupDb(ctx)
	if err != nil {
		return nil, err
	}

	repo := NewUserRepo(db)

	return repo, nil
}

func cleanUserRepo(repo *UserRepo) {
	defer repo.db.Close()

	_, err := repo.DeleteAll(context.Background())
	if err != nil {
		log.Println("Error tearing down the database")
	}
}

func insertUserDummyData(t *testing.T, repo *UserRepo) {
	dummyUser := storage.User{
		Id:          "",
		Email:       "john@gmail.com",
		CreatedAt:   time.Time{},
		Role:        storage.AuthRoleNone,
		CogUsername: "what-ever-username123",
		CogSub:      "123452431325-31525123",
		CogName:     "uagaewt1234",
		Disabled:    false,
	}

	_, err := repo.InsertMany(context.Background(), storage.UserList{dummyUser})
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

func TestUserRepo_GetByUsername(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupUserRepo(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	defer cleanUserRepo(repo)

	insertUserDummyData(t, repo)

	user, err := repo.GetByUsername(ctx, "what-ever-username123")
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if user.Email != " john@gmail.com" {
		t.Errorf("failed asserting email, got %s", user.Email)
	}
	if user.CogUsername != "what-ever-username123" {
		t.Errorf("failed asserting CogUsername, got %s", user.CogUsername)
	}
	if user.CogSub != "123452431325-31525123" {
		t.Errorf("failed asserting CogSub, got %s", user.CogSub)
	}
	if user.CogName != "uagaewt1234" {
		t.Errorf("failed asserting CogName, got %s", user.CogName)
	}
	if user.Role != storage.AuthRoleNone {
		t.Errorf("failed asserting Role, got %s", user.Role)
	}
}

func TestUserRepo_GetByUsername_NotFound(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()
	repo, err := setupUserRepo(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	defer cleanUserRepo(repo)

	insertUserDummyData(t, repo)

	_, err = repo.GetByUsername(ctx, "some-other")
	if !errors.As(err, &storage.NotFound{}) {
		t.Fatal("expected error of type not found")
	}
}

func TestUserRepo_Create(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupUserRepo(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	defer cleanUserRepo(repo)

	newUser := storage.UserCreationDto{
		Email:       "mary@gmail.com",
		Role:        storage.AuthRoleAdmin,
		CogUsername: "whatever-username",
		CogSub:      "whatever-sub",
		CogName:     "whatever-name",
		Disabled:    false,
	}

	user, err := repo.Create(ctx, newUser)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if user.Email != "mary@gmail.com" {
		t.Errorf("failed email assertion, got %s", user.Email)
	}
	if user.Role != storage.AuthRoleAdmin {
		t.Errorf("failed role assertion, got %s", user.Role)
	}
	if user.CogUsername != "whatever-username" {
		t.Errorf("failed CogUsername assertion, got %s", user.CogUsername)
	}
	if user.CogSub != "whatever-sub" {
		t.Errorf("failed CogSub assertion, got %s", user.CogSub)
	}
	if user.CogName != "whatever-name" {
		t.Errorf("failed CogName assertion, got %s", user.CogName)
	}
	if user.Disabled {
		t.Errorf("failed disabled assertion, got %s", user.CogName)
	}
}

func TestUserRepo_Create_Exists(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupUserRepo(ctx)
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
	defer cleanUserRepo(repo)

	insertUserDummyData(t, repo)

	newUser := storage.UserCreationDto{
		Email:       "john@gmail.com",
		Role:        storage.AuthRoleAdmin,
		CogUsername: "whatever-username",
		CogSub:      "whatever-sub",
		CogName:     "whatever-name",
		Disabled:    false,
	}

	_, err = repo.Create(ctx, newUser)
	if !errors.Is(err, storage.ErrDuplicate) {
		t.Errorf("expected error duplicate, got %v", err)
	}
}

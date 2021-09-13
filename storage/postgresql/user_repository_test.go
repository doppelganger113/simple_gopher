package postgresql

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"simple_gopher/storage"
	"simple_gopher/test"
	"testing"
	"time"
)

func setupUserRepo(ctx context.Context) (UserRepo, error) {
	db, err := setupDb(ctx)
	if err != nil {
		return UserRepo{}, err
	}

	repo := NewUserRepo(db)

	return repo, nil
}

func cleanUserRepo(repo UserRepo) {
	defer repo.db.Close()

	_, err := repo.DeleteAll(context.Background())
	if err != nil {
		log.Println("Error tearing down the database")
	}
}

func insertUserDummyData(t *testing.T, repo UserRepo) {
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
	assert.Nil(t, err)
}

func TestUserRepo_GetByUsername(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupUserRepo(ctx)
	assert.Nil(t, err)
	defer cleanUserRepo(repo)

	insertUserDummyData(t, repo)

	user, err := repo.GetByUsername(ctx, "what-ever-username123")
	assert.Nil(t, err)

	assert.Equal(t, "john@gmail.com", user.Email)
	assert.Equal(t, "what-ever-username123", user.CogUsername)
	assert.Equal(t, "what-ever-username123", user.CogUsername)
	assert.Equal(t, "123452431325-31525123", user.CogSub)
	assert.Equal(t, "uagaewt1234", user.CogName)
	assert.Equal(t, storage.AuthRoleNone, user.Role)
}

func TestUserRepo_GetByUsername_NotFound(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()
	repo, err := setupUserRepo(ctx)
	assert.Nil(t, err)
	defer cleanUserRepo(repo)

	insertUserDummyData(t, repo)

	_, err = repo.GetByUsername(ctx, "some-other")
	assert.Equal(t, true, errors.As(err, &storage.NotFound{}))
}

func TestUserRepo_Create(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupUserRepo(ctx)
	assert.Nil(t, err)
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
	assert.Nil(t, err)

	assert.Equal(t, "mary@gmail.com", user.Email)
	assert.Equal(t, storage.AuthRoleAdmin, user.Role)
	assert.Equal(t, "whatever-username", user.CogUsername)
	assert.Equal(t, "whatever-sub", user.CogSub)
	assert.Equal(t, "whatever-name", user.CogName)
	assert.Equal(t, false, user.Disabled)
}

func TestUserRepo_Create_Exists(t *testing.T) {
	test.SkipIfNotIntegrationTesting(t)
	ctx := context.Background()

	repo, err := setupUserRepo(ctx)
	assert.Nil(t, err)
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
	assert.Equal(t, true, errors.Is(err, storage.DuplicateErr))
}

package postgresql

import (
	"api/storage"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"strings"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(db *Database) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) GetByUsername(ctx context.Context, username string) (storage.User, error) {
	query := `SELECT id, email, role, cog_username, cog_sub, cog_name, created_at, updated_at, disabled
FROM users WHERE cog_username=$1 LIMIT 1 
`
	var user storage.User
	err := repo.db.dbPool.QueryRow(ctx, query, username).
		Scan(
			&user.Id,
			&user.Email,
			&user.Role,
			&user.CogUsername,
			&user.CogSub,
			&user.CogName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Disabled,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.User{}, storage.NotFound{Msg: "User not found by username " + username}
		}
		return storage.User{}, err
	}

	return user, nil
}

func (repo *UserRepo) Create(ctx context.Context, dto storage.UserCreationDto) (storage.User, error) {
	query := `INSERT INTO users
("email", "role", "cog_username", "cog_sub", "cog_name", "disabled")
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, role, cog_username, cog_sub, cog_name, created_at, updated_at, disabled 
`
	var user storage.User

	err := repo.db.dbPool.QueryRow(
		ctx,
		query,
		dto.Email,
		dto.Role,
		dto.CogUsername,
		dto.CogSub,
		dto.CogName,
		dto.Disabled,
	).
		Scan(
			&user.Id,
			&user.Email,
			&user.Role,
			&user.CogUsername,
			&user.CogSub,
			&user.CogName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Disabled,
		)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return storage.User{}, storage.ErrDuplicate
		}
		return storage.User{}, err
	}

	return user, nil
}

func (repo *UserRepo) InsertMany(ctx context.Context, users storage.UserList) (count int64, err error) {
	count, err = repo.db.dbPool.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"email", "role", "cog_username", "cog_sub", "cog_name", "disabled"},
		pgx.CopyFromSlice(len(users), func(i int) ([]interface{}, error) {
			return []interface{}{
				users[i].Email,
				string(users[i].Role),
				users[i].CogUsername,
				users[i].CogSub,
				users[i].CogName,
				users[i].Disabled,
			}, nil
		}),
	)
	if err != nil {
		return 0, err
	}

	return
}

func (repo *UserRepo) DeleteAll(ctx context.Context) (rowsAffected int64, err error) {
	query := "DELETE FROM users"
	cmdTag, err := repo.db.dbPool.Exec(ctx, query)
	if err != nil {
		return 0, err
	}

	rowsAffected = cmdTag.RowsAffected()
	return
}

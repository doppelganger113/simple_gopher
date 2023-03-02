package postgresql

import (
	"api/logger"
	"context"
	"errors"
	"os"
)

func setupDb(ctx context.Context) (*Database, error) {
	dbUrl := os.Getenv("DATABASE_TEST_URL")
	if dbUrl == "" {
		return nil, errors.New("missing DATABASE_TEST_URL env variable")
	}

	db := NewDatabase(logger.NewLogger())
	err := db.Connect(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	return db, nil
}

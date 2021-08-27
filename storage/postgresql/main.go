package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"simple_gopher/cloud_patterns"
	"time"
)

type Database struct {
	dbPool *pgxpool.Pool
}

func NewDatabase() *Database {
	return &Database{}
}

// Connect method connects the database where the connection url is in format:
// "postgresql://postgres:example@localhost/dbname"
func (db *Database) Connect(ctx context.Context, connectionString string) error {
	log.Info().Msg("Trying to connect to the database")

	retry := cloud_patterns.NewRetry(3 * time.Second)
	err := retry.Execute(ctx, func(ctx context.Context, retryCount uint) error {
		if retryCount > 0 {
			log.Info().Msgf("Retrying connection %d time", retryCount)
		}

		dbpool, err := pgxpool.Connect(context.Background(), connectionString)
		if err != nil {
			return err
		}
		db.dbPool = dbpool

		return nil
	})
	if err != nil {
		return errors.New("failed to connect")
	}

	log.Info().Msg("Connected to the database")

	return nil
}

func (db *Database) Close() {
	log.Info().Msg("Closing database connection.")
	db.dbPool.Close()
}

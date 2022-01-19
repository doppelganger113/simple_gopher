package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"simple_gopher/concurrency"
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
	log.Info().Msg("[Database]: Trying to connect...")

	retry := concurrency.NewRetry(3, 3*time.Second)
	err := retry.Execute(ctx, func(ctx context.Context, retryCount uint) error {
		if retryCount > 0 {
			log.Info().Msgf("[Database]: Retrying connection %d", retryCount)
		}

		dbpool, err := pgxpool.Connect(context.Background(), connectionString)
		if err != nil {
			return err
		}
		db.dbPool = dbpool

		return nil
	})
	if err != nil {
		return errors.New("[Database]: failed to connect")
	}

	log.Info().Msg("[Database]: connected")

	return nil
}

func (db *Database) Close() {
	log.Info().Msg("[Database]: Closing connection.")
	db.dbPool.Close()
}

package postgresql

import (
	"api/pkg/concurrency"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"time"
)

const connectionRetryCount = 4
const connectionRetryBackoff = 3 * time.Second

type Database struct {
	dbPool *pgxpool.Pool
	logger *zerolog.Logger
}

func NewDatabase(logger *zerolog.Logger) *Database {
	return &Database{logger: logger}
}

// Connect method connects the database where the connection url is in format:
// "postgresql://postgres:example@localhost/dbname"
func (db *Database) Connect(ctx context.Context, connectionString string) error {
	db.logger.Info().Msg("[Database]: Trying to connect...")

	retry := concurrency.NewRetry(connectionRetryCount, connectionRetryBackoff)
	err := retry.Execute(ctx, func(ctx context.Context, retryCount uint) error {
		if retryCount > 0 {
			db.logger.Info().Msg("[Database]: Retrying connection...")
		}

		dbpool, err := pgxpool.Connect(context.Background(), connectionString)
		if err != nil {
			return err
		}
		db.dbPool = dbpool

		return nil
	})
	if err != nil {
		return err
	}

	db.logger.Info().Msg("[Database]: connected")
	return nil
}

func (db *Database) Close() {
	db.logger.Info().Msg("[Database]: Closing connection.")
	db.dbPool.Close()
}

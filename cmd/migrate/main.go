package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"time"
)

const (
	defaultMigrationDir = "resources/migrations"
)

type settings struct {
	steps            int
	force            int
	connectionUrl    string
	migrationDirPath string
}

func main() {
	sett, err := getSettings()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connecting to the database")
	db, err := sql.Open("postgres", sett.connectionUrl)
	if err != nil {
		log.Fatalf("failed connecting to postgres: %v\n", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln(fmt.Errorf("failed creating DB driver: %w", err))
	}

	log.Println("Reading migration from: " + sett.migrationDirPath)
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", sett.migrationDirPath), "postgres", driver,
	)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed creating new db instance: %w", err))
	}

	// Timer
	start := time.Now()
	defer func() {
		end := time.Now()
		log.Printf("Migration completed after: %s\n", end.Sub(start).String())
	}()

	if sett.force != 0 {
		err = m.Force(sett.force)
	} else if sett.steps != 0 {
		log.Println(getMigrationDescription(sett.steps))
		err = m.Steps(sett.steps)
	} else {
		log.Println("Migrating all the way up")
		err = m.Up()
	}

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No changes detected for migration.")
		} else {
			log.Fatalln(fmt.Errorf("migration failed: %w", err))
		}
	}
}

func getMigrationDescription(steps int) string {
	if steps < 0 {
		return fmt.Sprintf("Performing migration down %d steps\n", steps)
	}

	return fmt.Sprintf("Performing migration up %d steps\n", steps)
}

func getSettings() (*settings, error) {
	// Optional
	steps := flag.Int("steps", 0, "if steps > 0 migrate up, if steps < 0 migrate down")
	force := flag.Int("f", 0, "force migration to specified version")
	migrationDirPath := flag.String(
		"dir", defaultMigrationDir, "Migration directory path",
	)

	flag.Parse()
	connectionUrl := os.Getenv("DATABASE_URL")
	if connectionUrl == "" {
		return nil, errors.New("missing 'DATABASE_URL' env variable")
	}

	return &settings{
		steps:            *steps,
		force:            *force,
		connectionUrl:    connectionUrl,
		migrationDirPath: *migrationDirPath,
	}, nil
}

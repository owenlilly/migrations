package migrations

import (
	"errors"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	MigrateUp   Direction = "up"
	MigrateDown Direction = "down"
)

type Direction string

type Config struct {
	Direction   Direction // up or down
	Steps       int       // number of migration steps, runs all if < 0, runs n if > 0, no change if 0
	DatabaseURL string    // database connection string
	SourceURL   string    // directory or url containing migration files
}

func RunMigrations(p Config) error {
	m, err := migrate.New(p.SourceURL, p.DatabaseURL)
	if err != nil {
		return err
	}

	if p.Direction == MigrateUp {
		err = runUp(m, p.Steps)
	} else if p.Direction == MigrateDown {
		err = runDown(m, p.Steps)
	}

	if err != nil && err.Error() == "no change" {
		return nil
	}

	return err
}

func ResetAllData(sourceURL, connectionString string) error {
	if !strings.Contains(connectionString, "_test") &&
		!strings.Contains(connectionString, "_demo") {
		return errors.New("cannot reset non-test database")
	}

	if !strings.Contains(connectionString, "_test") &&
		!strings.Contains(connectionString, "_demo") {
		return errors.New("cannot reset non-test database")
	}
	m, err := migrate.New("file://"+sourceURL, connectionString)
	if err != nil {
		return err
	}

	return m.Drop()
}

func runUp(m *migrate.Migrate, upSteps int) (err error) {
	if upSteps == 0 {
		// no change
		return
	}

	if upSteps < 0 {
		// all the way up
		return m.Up()
	}

	// n step up
	return m.Steps(upSteps)
}

func runDown(m *migrate.Migrate, downSteps int) (err error) {
	if downSteps == 0 {
		// no change
		return
	}

	if downSteps < 0 {
		// all the way downSteps
		return m.Down()
	}

	// n steps downSteps
	return m.Steps(-downSteps)
}

package database

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// 🔥 Embed migration SQL files
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations executes all embedded migrations
func migrateDb(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}

	// ✅ Use embedded migrations instead of filesystem
	d, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return errors.Wrap(err, "loading embedded migrations")
	}

	migrator, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "creating migrator")
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "executing migration")
	}

	currentVersion, dirty, err := migrator.Version()
	if err != nil {
		return errors.Wrap(err, "getting migration version")
	}

	logrus.WithFields(logrus.Fields{
		"version": currentVersion,
		"dirty":   dirty,
	}).Debug("Database migrated")

	return nil
}

// func migrateDb(db *sql.DB) error {
// 	driver, err := postgres.WithInstance(db, &postgres.Config{})
// 	if err != nil {
// 		return errors.Wrap(err, "connecting to database")
// 	}

// 	migrationSource := fmt.Sprintf("file://%sinternal/database/migrations/", *config.DataDirectory)
// 	migrator, err := migrate.NewWithDatabaseInstance(migrationSource, "postgres", driver)
// 	if err != nil {
// 		return errors.Wrap(err, "creating migrator")
// 	}

// 	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
// 		return errors.Wrap(err, "executing migration")
// 	}

// 	currentVersion, dirty, err := migrator.Version()
// 	if err != nil {
// 		return errors.Wrap(err, "getting migration version")
// 	}

// 	logrus.WithFields(logrus.Fields{
// 		"version": currentVersion,
// 		"dirty":   dirty,
// 	}).Debug("Database migrated")

// 	return nil
// }

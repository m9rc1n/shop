package repository

import (
	"fmt"
	migrater "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/m9rc1n/shop/pkg/log"
	"os"
)

// Migrator automatically migrates the database
type Migrator interface {
	// MigrateUp item
	MigrateUp() error
	// MigrateDown item
	MigrateDown() error
}

type migrator struct {
	db     *sqlx.DB
	logger log.Logger
}

// NewMigrator creates a new migrator
func NewMigrator(db *sqlx.DB, log log.Logger) Migrator {
	return migrator{db: db, logger: log}
}

// Migrate to run database migration up or down
func (r migrator) configureMigrater() (*migrater.Migrate, error) {
	path, err := os.Getwd()
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	migrationPath := fmt.Sprintf("file://%s/migration", path)
	r.logger.Infof("migrationPath : %s", migrationPath)
	driver, err := postgres.WithInstance(r.db.DB, &postgres.Config{})
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return migrater.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
}

// MigrateUp item
func (r migrator) MigrateUp() error {
	m, err := r.configureMigrater()
	if err != nil {
		r.logger.Errorf("An error occurred while configuring database migrater: %v", err)
		return err
	}
	r.logger.Info("Migrate up")
	if err := m.Up(); err != nil && err != migrater.ErrNoChange {
		r.logger.Errorf("An error occurred while syncing the database: %v", err)
		return err
	}
	return nil
}

// MigrateDown item
func (r migrator) MigrateDown() error {
	m, err := r.configureMigrater()
	if err != nil {
		r.logger.Errorf("An error occurred while configuring database migrater: %v", err)
		return err
	}
	r.logger.Info("Migrate down")
	if err := m.Down(); err != nil && err != migrater.ErrNoChange {
		r.logger.Errorf("An error occurred while syncing the database: %v", err)
		return err
	}
	return nil
}

package repository

import (
	"fmt"
	migrater "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"os"
)

// Migrate to run database migration up or down
func (r repository) configureMigrater() (*migrater.Migrate, error) {
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
func (r repository) MigrateUp() error {
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
func (r repository) MigrateDown() error {
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

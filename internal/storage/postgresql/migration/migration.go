package migration

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func MigrationsUp(db *sqlx.DB, dsn string) error {
	const op = "internal.storage.postgresql.migration.MigrationsUp()"

	if db == nil {
		return fmt.Errorf("%s:%s", op, "database isn`t established")
	}

	migrationDB, err := connectionForMigration(dsn)
	if err != nil {
		return err
	}

	sqlDatabase := migrationDB.DB
	driver, err := newMigrationDriver(sqlDatabase)
	if err != nil {
		return err
	}

	defer closeMigration(driver, migrationDB, op)

	migration, err := newMigrationInstance(driver)
	if err != nil {
		return err
	}

	err = startMigrationUp(migration)
	if err != nil {
		return err
	}

	return nil
}

func MigrationsDown(db *sqlx.DB, dsn string) error {
	const op = "internal.storage.postgresql.migration.MigrationsDown()"

	if db == nil {
		return fmt.Errorf("%s:%s", op, "database isn`t established")
	}

	migrationDB, err := connectionForMigration(dsn)
	if err != nil {
		return err
	}

	sqlDatabase := migrationDB.DB
	driver, err := newMigrationDriver(sqlDatabase)
	if err != nil {
		return err
	}

	defer closeMigration(driver, migrationDB, op)

	migration, err := newMigrationInstance(driver)
	if err != nil {
		return err
	}

	err = startMigrationDown(migration)
	if err != nil {
		return err
	}

	return nil
}

func connectionForMigration(dsn string) (*sqlx.DB, error) {
	const op = "internal.storage.postgresql.migration.connectionForMigration()"
	migration, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s:%w(%s)", op, err, "failed to create connection for migrations")
	}
	return migration, nil
}

func newMigrationDriver(db *sql.DB) (database.Driver, error) {
	const op = "internal.storage.postgresql.migration.newMigrationDriver()"
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s:%w(%s)", op, err, "couldn't create driver")
	}
	return driver, nil
}

func closeMigration(driver database.Driver, migrationDB *sqlx.DB, op string) {
	if err := driver.Close(); err != nil {
		log.Printf("%s:%s", op, "migration's driver couldn't close")
	}

	if err := migrationDB.Close(); err != nil {
		log.Printf("%s:%s", op, "migration's connection couldn't close")
	}
}

func newMigrationInstance(driver database.Driver) (*migrate.Migrate, error) {
	const op = "internal.storage.postgresql.migration.newMigrationInstance()"
	migrationExmpl, err := migrate.NewWithDatabaseInstance(
		"file://internal/storage/postgresql/migration",
		"postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("%s:%w(%s)", op, err, "coudn't create migrate instance")
	}
	return migrationExmpl, nil
}

func startMigrationUp(migration *migrate.Migrate) error {
	const op = "internal.storage.postgresql.migration.startMigrationUp"
	err := migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s:%w(%s)", op, err, "failed to apply migrations")
	}

	if err == migrate.ErrNoChange {
		log.Printf("%s:%s", op, "no migrations to apply")
	} else {
		log.Printf("%s:%s", op, "migrations applied successfully")
	}
	return nil
}

func startMigrationDown(migration *migrate.Migrate) error {
	const op = "internal.storage.postgresql.migration.startMigrationUp"
	err := migration.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s:%w(%s)", op, err, "failed to apply migrations")
	}

	if err == migrate.ErrNoChange {
		log.Printf("%s:%s", op, "no migrations to apply")
	} else {
		log.Printf("%s:%s", op, "migrations applied successfully")
	}
	return nil
}

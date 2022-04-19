package migrate

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbUrl string, migrationsPath string) error {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return err
	}
	return m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
}

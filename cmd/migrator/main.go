package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	var migrationsPath, migrationsTable, host, user, password, dbname string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.StringVar(&host, "host", "", "PostgreSQL host")
	flag.StringVar(&user, "user", "", "PostgreSQL user")
	flag.StringVar(&password, "password", "", "PostgreSQL password")
	flag.StringVar(&dbname, "dbname", "", "PostgreSQL database name")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	if host == "" || user == "" || password == "" || dbname == "" {
		panic("PostgreSQL connection info is required")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}
		panic(err)
	}

	fmt.Println("migrations applied")
}

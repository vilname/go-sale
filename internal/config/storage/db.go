package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"os"
)

var db *pgxpool.Pool

func InitDB() {

	// не получилось написать миграцию для библиотеки pgx по этому сперва создаю подключение
	// через sql.DB и запускаю миграции, а тотом создаю подключение pgx
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbSql, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dbHost, dbUser, dbPass, dbName, dbPort),
	)

	driver, _ := postgres.WithInstance(dbSql, &postgres.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://migration",
		"postgres", driver)

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		_ = fmt.Errorf("error when migration up: %v", err)
	}

	dbSql.Close()

	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func GetDB() *pgxpool.Pool {
	return db
}

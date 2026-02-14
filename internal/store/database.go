package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

//! getEnv --> helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

//! Open --> establishes connection to PostgreSQL database
//! Using port 5445 locally (not 5432) to avoid Windows port reservation conflicts
//! In Docker, it uses environment variables and connects to port 5432
func Open() (*sql.DB, error) {
	//* get database configuration from environment variables with fallbacks
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5445") //* 5445 for local Windows, 5432 in Docker
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "postgres")

	//* connection string with all database credentials
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := sql.Open("pgx", connStr)

	//? if caught any error while opening connection
	if err != nil {
		return nil, fmt.Errorf("db : open %w", err)
	}
	fmt.Printf("Connected to the Database at %s:%s...\n", host, port)
	return db, err //* return connection pool

}

//! Migratefs --> runs database migrations from embedded filesystem
//! Migrations are version control for database schema changes
func Migratefs(db *sql.DB,migrationfs fs.FS,dir string) error {
	//* tell goose to read migrations from embedded FS (not disk)
	goose.SetBaseFS(migrationfs)

	defer func() {
		goose.SetBaseFS(nil) //* cleanup after migrations run
	}()
	
	return Migrate(db,dir) //* execute migrations

}
//! Migrate --> applies pending database migrations using goose
func Migrate(db *sql.DB,dir string)error{
	//* set database dialect to postgres
	err := goose.SetDialect("postgres")

	if err != nil {
		return fmt.Errorf("migrate : %w ",err)
	}
	//* run goose up to apply all pending migrations
	err = goose.Up(db,dir)
	if err != nil {
		return fmt.Errorf("goose up : %w ",err)
	}
	return nil

}
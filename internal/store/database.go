package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

//! Open --> establishes connection to PostgreSQL database
//! Using port 5445 (not 5432) to avoid Windows port reservation conflicts
func Open() (*sql.DB, error) {
	//* connection string with all database credentials
	db,err := sql.Open("pgx","host=localhost user=postgres password=postgres dbname=postgres port=5445 sslmode=disable")

	//? if caught any error while opening connection
	if err != nil {
		return nil,fmt.Errorf("db : open %w", err)
	}
	fmt.Println("Connected to the Database...")
	return db,err //* return connection pool

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
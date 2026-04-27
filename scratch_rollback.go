package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"grc_be/internal/data/schema"
)

func main() {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=Tambunan140705 dbname=grc_db sslmode=disable")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	goose.SetBaseFS(schema.Migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	fmt.Println("WIPING DATABASE SCHEMA (dropping public schema)...")
	db.Exec("DROP SCHEMA public CASCADE;")
	db.Exec("CREATE SCHEMA public;")
	db.Exec("GRANT ALL ON SCHEMA public TO postgres;")
	db.Exec("GRANT ALL ON SCHEMA public TO public;")

	fmt.Println("Database wiped clean. Now run the app to reconstruct everything.")
}

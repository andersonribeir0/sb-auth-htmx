package main

import (
	"embed"
	"flag"
	"fmt"
	"log"

	"dreampicai/internal/database"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	var migrate string
	flag.StringVar(&migrate, "migrate", "", "Direction to migrate the database (up or down)")
	flag.Parse()

	db := database.NewMigrationSvcProvider()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("pgx"); err != nil {
		panic(err)
	}

	switch migrate {
	case "up":
		if err := goose.Up(db.DB().DB, "migrations"); err != nil {
			panic(err)
		}
	case "down":
		if err := goose.Down(db.DB().DB, "migrations"); err != nil {
			panic(err)
		}
	default:
		fmt.Println("Invalid migrate flag value. Use 'up' or 'down'.")
	}
}

package main

import (
	"context"
	"log"

	"dreampicai/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewMigrationSvcProvider()

	tables := []string{
		"accounts",
		"goose_db_version",
	}

	_, err = db.DB().NewDropTable().Table(tables...).Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

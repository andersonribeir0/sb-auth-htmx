package main

import (
	"context"

	"dreampicai/internal/database"
)

func main() {
	db := database.NewMigrationSvcProvider()

	tables := []string{
		"accounts",
		"goose_db_version",
	}

	_, err := db.DB().NewDropTable().Table(tables...).Exec(context.Background())
	if err != nil {
		panic(err)
	}
}

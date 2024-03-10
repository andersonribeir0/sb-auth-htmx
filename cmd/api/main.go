package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"dreampicai/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	server := handler.NewServer()

	slog.Info("application running", "port", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

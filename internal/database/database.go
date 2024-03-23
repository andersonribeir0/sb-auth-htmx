package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Service interface {
	Health() map[string]string
	DB() *sql.DB
}

type service struct {
	db *bun.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", username, password, host, port, database)
	slog.Info("connStr", "value", connStr)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		slog.Error("sql open err", "err", err)
		log.Fatal(err)
	}

	s := &service{db: bun.NewDB(db, pgdialect.New())}

	if len(os.Getenv("DEBUG")) > 0 {
		s.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return s
}

func (s *service) DB() *sql.DB {
	return s.db.DB
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

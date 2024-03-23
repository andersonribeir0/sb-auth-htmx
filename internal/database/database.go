package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"dreampicai/types"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Service interface {
	Health() map[string]string
	CreateAccount(ctx context.Context, account *types.Account) error
}

type MigrationServiceProvider interface {
	DB() *bun.DB
}

type service struct {
	db *bun.DB
}

func NewMigrationSvcProvider() MigrationServiceProvider {
	return newSvc()
}

func New() Service {
	return newSvc()
}

func newSvc() *service {
	var (
		host     = os.Getenv("DB_HOSTNAME")
		database = os.Getenv("DB_DATABASE")
		password = os.Getenv("DB_PASSWORD")
		username = os.Getenv("DB_USERNAME")
		port     = os.Getenv("DB_PORT")
	)
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s", host, username, password, port, database)
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

func (s *service) DB() *bun.DB {
	return s.db
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

func (s *service) CreateAccount(ctx context.Context, account *types.Account) error {
	_, err := s.db.NewInsert().Model(account).Exec(ctx)
	return err
}

package types

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int `bun:"id,pk,autoincrement"`
	UserID    uuid.UUID
	Username  string
	CreatedAt time.Time `bun:"default:'now()'"`
}

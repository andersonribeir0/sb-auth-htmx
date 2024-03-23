package types

import "github.com/google/uuid"

const (
	UserContextKey = "user"
	AccessTokenKey = "accessToken"
)

type AuthenticatedUser struct {
	ID         uuid.UUID
	Email      string
	IsLoggedIn bool
}

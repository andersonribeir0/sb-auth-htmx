package types

const (
	UserContextKey = "user"
	AccessTokenKey = "accessToken"
)

type AuthenticatedUser struct {
	Email      string
	IsLoggedIn bool
}

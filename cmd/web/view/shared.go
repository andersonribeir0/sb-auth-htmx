package view

import (
	"context"
	"dreampicai/types"
	"log/slog"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserContextKey).(types.AuthenticatedUser)
	slog.Info("user", "email", user.Email)
	if !ok {
		return types.AuthenticatedUser{}
	}

	return user
}

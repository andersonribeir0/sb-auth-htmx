package handler

import (
	"log/slog"
	"net/http"

	"dreampicai/cmd/web/view/home"
	"dreampicai/types"
)

func (s *Server) HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	account := &types.Account{
		UserID:   user.ID,
		Username: "foobarbaz",
	}

	if err := s.db.CreateAccount(r.Context(), account); err != nil {
		slog.Error("creating account", "err", err)
		return err
	}

	slog.Info("user created", "account", account)

	return home.Index().Render(r.Context(), w)
}

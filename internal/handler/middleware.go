package handler

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"dreampicai/internal/database"
	"dreampicai/pkg/sb"
	"dreampicai/types"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

func RedirectIfAccountExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := getAuthenticatedUser(r)
		if user.IsLoggedIn {
			account, err := database.GetInstance().GetAccountByUserID(r.Context(), user.ID.String())
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				slog.Error("Error fetching account", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if len(account.UserID.String()) > 0 {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func WithAccount(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := getAuthenticatedUser(r)
		account, err := database.GetInstance().GetAccountByUserID(r.Context(), user.ID.String())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/account/setup", http.StatusSeeOther)
				return
			}
			const errMsg = "could not fetch account data"
			slog.Error(errMsg, "err", err)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
		user.Account = account
		slog.Info("account", "data", user)

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func WithAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		user := getAuthenticatedUser(r)
		if !user.IsLoggedIn {
			hxRedirect(w, r, "/login")
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/public") {
			next.ServeHTTP(w, r)
			return
		}

		store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
		sess, err := store.Get(r, types.UserContextKey)
		if err != nil || len(sess.Values) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		accessToken, ok := sess.Values[types.AccessTokenKey]
		if !ok {
			next.ServeHTTP(w, r)
		}

		resp, err := sb.Client.Auth.User(r.Context(), accessToken.(string))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		user := types.AuthenticatedUser{
			ID:         uuid.MustParse(resp.ID),
			Email:      resp.Email,
			IsLoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

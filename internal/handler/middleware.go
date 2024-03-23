package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"dreampicai/pkg/sb"
	"dreampicai/types"

	"github.com/gorilla/sessions"
)

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
			Email:      resp.Email,
			IsLoggedIn: true,
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

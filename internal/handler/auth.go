package handler

import (
	"log/slog"
	"net/http"

	"dreampicai/cmd/web/view/auth"
	"dreampicai/pkg/sb"
	"dreampicai/pkg/util"

	"github.com/nedpals/supabase-go"
)

func (s *Server) HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func (s *Server) HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if !util.ValidateEmail(credentials.Email) {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{Email: "Please enter a valid email."}))
	}

	if !util.ValidateStrongPassword(credentials.Password) {
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{Password: "Password must be 8+ characters, including uppercase, lowercase, numbers, and symbols."}))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "Invalid credentials.",
		}))
	}

	cookie := &http.Cookie{
		Value:    resp.AccessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

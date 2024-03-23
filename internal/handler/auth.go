package handler

import (
	"log/slog"
	"net/http"
	"os"

	"dreampicai/cmd/web/view/auth"
	"dreampicai/pkg/kit/validate"
	"dreampicai/pkg/sb"
	"dreampicai/types"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

func (s *Server) HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Signup())
}

func (s *Server) HandleSignupPost(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	errors := auth.SignupErrors{}

	if ok := validate.New(&params, validate.Fields{
		"Email":           validate.Rules(validate.Email, validate.Required),
		"Password":        validate.Rules(validate.Password, validate.Required),
		"ConfirmPassword": validate.Rules(validate.Equal(params.Password), validate.Message("Passwords must match.")),
	}).Validate(&errors); !ok {
		return render(r, w, auth.SignupForm(params, errors))
	}

	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		slog.Error("signup error", "err", err)
		return render(r, w, auth.SignupForm(params, auth.SignupErrors{SignupErr: "Signup failed."}))
	}

	slog.Info("user", "data", user)

	return render(r, w, auth.SignupSuccess(user.Email))
}

func (s *Server) HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func (s *Server) HandleLoginPost(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	var errors auth.LoginErrors
	if ok := validate.New(&credentials, validate.Fields{
		"Email":    validate.Rules(validate.Email, validate.Required),
		"Password": validate.Rules(validate.Password, validate.Required),
	}).Validate(&errors); !ok {
		return render(r, w, auth.LoginForm(credentials, errors))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "Invalid credentials.",
		}))
	}

	setAuthCookie(w, r, resp.AccessToken)

	return hxRedirect(w, r, "/")
}

func (s *Server) HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}

	setAuthCookie(w, r, accessToken)

	return hxRedirect(w, r, "/")
}

func (s *Server) HandleLogoutPost(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	sess, _ := store.Get(r, types.UserContextKey)
	sess.Values[types.AccessTokenKey] = ""
	if err := sess.Save(r, w); err != nil {
		return err
	}

	return hxRedirect(w, r, "/")
}

func (s *Server) HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: os.Getenv("GOOGLE_LOGIN_CALLBACK_URL"),
	})
	if err != nil {
		return err
	}

	return hxRedirect(w, r, resp.URL)
}

func setAuthCookie(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	sess, _ := store.Get(r, types.UserContextKey)
	sess.Values[types.AccessTokenKey] = accessToken
	if err := sess.Save(r, w); err != nil {
		return err
	}

	return sess.Save(r, w)
}

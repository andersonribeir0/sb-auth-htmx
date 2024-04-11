package handler

import (
	"log/slog"
	"net/http"

	"dreampicai/cmd/web/view/settings"
	"dreampicai/pkg/kit/validate"
)

func (s *Server) HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	return render(r, w, settings.Index(user))
}

func (s *Server) HandleUpdateProfilePut(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)

	var (
		params = settings.ProfileParams{Username: r.FormValue("username")}
		errors settings.ProfileErrors
	)
	if ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(3), validate.Max(50)),
	}).Validate(&errors); !ok {
		return render(r, w, settings.ProfileForm(params, errors))
	}

	user.Account.Username = params.Username

	err := s.db.UpdateUsername(r.Context(), &user.Account)
	if err != nil {
		slog.Error("failed to update username", "username", params.Username)
		return err
	}

	return render(r, w, settings.ProfileForm(params, settings.ProfileErrors{}))
}
